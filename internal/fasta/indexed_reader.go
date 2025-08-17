package fasta

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/guillechuma/bio-tui/internal/index"
)

// IndexedFastaReader manages access to a FASTA file using a .fai index.
type IndexedReader struct {
	file  *os.File                   // The open FASTA file handle
	Index map[string]index.FaiRecord // The in-memory index, mapping sequence ID to its record
}

// NewIndexedReader creates a reader by opening a FASTA file and parsing its .fai index.
func NewIndexedReader(fastaPath string) (*IndexedReader, error) {
	indexPath := fastaPath + ".fai"

	// Check if the index file exists.
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// if it does not exists, build it
		fmt.Fprintf(os.Stderr, "[info] FASTA index not found. Building now at %s...\n", indexPath)
		if err := index.BuildFai(fastaPath); err != nil {
			return nil, fmt.Errorf("failed to build FASTA index: %w", err)
		}
	}
	// 1. Open the main FASTA file. Keep it open.
	fastaFile, err := os.Open(fastaPath)
	if err != nil {
		return nil, fmt.Errorf("could not open fasta file: %w", err)
	}

	// 2. Open and parse the index file with ParseFai
	idx, err := index.ParseFai(indexPath)
	if err != nil {
		fastaFile.Close() // Clean up the already opened fasta file
		return nil, err
	}
	// Create the reader instance and return it
	reader := &IndexedReader{
		file:  fastaFile,
		Index: idx,
	}
	return reader, nil
}

// Fetch retrieves a single FastaRecord by its ID.
func (r *IndexedReader) Fetch(id string) (*FastaRecord, error) {
	// 1. Look up the record in our in-memory index.
	indexRecord, ok := r.Index[id]
	if !ok {
		return nil, fmt.Errorf("sequence with id '%s' not found in index", id)
	}
	// 2. Seek to the start of the sequence data in the large FASTA file.
	_, err := r.file.Seek(indexRecord.Offset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to offset for id '%s': %w", id, err)
	}

	// 3. Read the sequence data directly, stripping newlines on the fly.
	// We know the total length, how many bases are on each line,
	// and how many bytes each line takes up (including the newline).
	var seqBuilder bytes.Buffer
	basesToRead := indexRecord.Length
	for basesToRead > 0 {
		// Determine how many bases to read from the current line.
		chunkSize := min(basesToRead, indexRecord.LineBases)

		// Read the full line (base + newline) from the file.
		lineBuffer := make([]byte, indexRecord.LineBytes)
		bytesRead, err := r.file.Read(lineBuffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read sequence data for id '%s': %w", id, err)
		}

		// Append only the bases (not the newline) to our sequence.
		seqBuilder.Write(lineBuffer[:chunkSize])

		// In case we hit EOF in the middle of a line.
		if bytesRead < int(indexRecord.LineBytes) && err == io.EOF {
			break
		}

		basesToRead -= chunkSize
	}

	// 4. Manually construct the FastaRecord.
	record := &FastaRecord{
		ID:  indexRecord.Name,
		Seq: seqBuilder.Bytes(),
	}
	record.Type = InferSequenceType(record.Seq) // Infer type as before

	return record, nil
}

// Close closes the underlying FASTA file.
func (r *IndexedReader) Close() error {
	return r.file.Close()
}
