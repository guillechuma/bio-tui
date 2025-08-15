package fasta

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Parser struct {
	scanner    *bufio.Scanner
	peekedLine string // Store the next header line we've already read
}

// NewParser creates a new FASTA parser
func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
	}
}

// Next returns the next FastaRecord from the stream.
// It returns io.EOF when the stream is exhausted.
func (p *Parser) Next() (*FastaRecord, error) {
	var headerLine string

	// Step 1: Find the header for the current record.
	// It might be peekedLine from last call, or we need to scan.
	if p.peekedLine != "" {
		headerLine = p.peekedLine
		p.peekedLine = "" // Clear it now that we're using it
	} else {
		for p.scanner.Scan() {
			line := p.scanner.Text()
			if strings.HasPrefix(line, ">") {
				headerLine = line
				break // Found it, stop scanning
			}
		}
	}

	// If we still don't have a header, we've reached the end of the file.
	if headerLine == "" {
		// Check for any scanner errors before done.
		if err := p.scanner.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF // Standard way to signal completion.
	}

	// Step 2: We have a header. Create the record and parse header.
	record := &FastaRecord{}
	parseHeaderLine(headerLine, record)

	// Step 3: Read sequence lines until the next header or EOF
	// A bytes.Buffer is the most efficient way yo build the seq string.
	var seqBuilder bytes.Buffer
	for p.scanner.Scan() {
		line := p.scanner.Text()
		if strings.HasPrefix(line, ">") {
			// Found the start of the *next* record.
			// Save it for the next call to Next() and stop here.
			p.peekedLine = line
			break
		}
		// If it's not a header, it's a sequence line.
		seqBuilder.WriteString(strings.TrimSpace(line))
	}

	record.Seq = seqBuilder.Bytes()

	return record, nil
}

// parseHeaderLine takes a header line (e.g., ">ID Description")
// and sets the fields in the FastaRecord.
func parseHeaderLine(line string, record *FastaRecord) {
	headerText := line[1:]
	parts := strings.SplitN(headerText, " ", 2)
	record.ID = parts[0]
	if len(parts) > 1 {
		record.Description = parts[1]
	}
}
