package index

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/guillechuma/bio-tui/internal/adapters/fasta"
)

// FaiRecord holds the index information for a single sequence in a FASTA file.
type FaiRecord struct {
	Name      string // Name of the sequence
	Length    int64  // Total length of the sequence in bases
	Offset    int64  // Byte offset in the file where the sequence starts
	LineBases int64  // Number of bases per line
	LineBytes int64  // Number of bytes per line (including newline)
}

// IndexedFastaReader manages access to a FASTA file using a .fai index.
type IndexedFastaReader struct {
	file  *os.File             // The open FASTA file handle
	index map[string]FaiRecord // The in-memory index, mapping sequence ID to its record
}

// NewIndexedFastaReader creates a reader for a FASTA file by loading its .fai index.
func NewIndexedFastaReader(fastaPath string) (*IndexedFastaReader, error) {
	// 1. Open the main FASTA file. Keep it open.
	fastaFile, err := os.Open(fastaPath)
	if err != nil {
		return nil, fmt.Errorf("could not open fasta file: %w", err)
	}

	// 2. Open and parse the index file.
	indexPath := fastaPath + ".fai"
	indexFile, err := os.Open(indexPath)
	if err != nil {
		fastaFile.Close() // Clean up the already opened fasta file
		return nil, fmt.Errorf("could not open index file (.fai): %w", err)
	}
	defer indexFile.Close()

	// 3. Read the index line-by-line and store it in a map.
	index := make(map[string]FaiRecord)
	scanner := bufio.NewScanner(indexFile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) != 5 {
			continue // Skip malformed lines
		}
		length, _ := strconv.ParseInt(parts[1], 10, 64)
		offset, _ := strconv.ParseInt(parts[2], 10, 64)
		lineBases, _ := strconv.ParseInt(parts[3], 10, 64)
		lineBtyes, _ := strconv.ParseInt(parts[4], 10, 64)

		record := FaiRecord{
			Name:      parts[0],
			Length:    length,
			Offset:    offset,
			LineBases: lineBases,
			LineBytes: lineBtyes,
		}
		index[record.Name] = record
	}

	// 4. Create the reader instance and return it
	reader := &IndexedFastaReader{
		file:  fastaFile,
		index: index,
	}
	return reader, nil
}

// Fetch retrieves a single FastaRecord by its ID.
func (r *IndexedFastaReader) Fetch(id string) (*fasta.FastaRecord, error) {
	// 1. Look up the record in our in-memory index.
	indexRecord, ok := r.index[id]
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
	record := &fasta.FastaRecord{
		ID:  indexRecord.Name,
		Seq: seqBuilder.Bytes(),
	}
	record.Type = fasta.InferSequenceType(record.Seq) // Infer type as before

	return record, nil
}

// Close closes the underlying FASTA file.
func (r *IndexedFastaReader) Close() error {
	return r.file.Close()
}
