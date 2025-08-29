package fastq

import (
	"bufio"
	"fmt"
	"io"
)

// Parser reads FastqRecords from a reader.
type Parser struct {
	scanner *bufio.Scanner
}

// NewParser creates a new FASTQ parser
func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
	}
}

// Next returns the next FastqRecord from the stream
func (p *Parser) Next() (*FastqRecord, error) {
	// FASTQ record is always four lines.

	// Read the first line (ID). If it fails, we might be at the end of the file.
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
	idLine := p.scanner.Text()

	// Read the next three lines (sequence, separator, quality).
	if !p.scanner.Scan() {
		return nil, fmt.Errorf("unexpected EOF after id line")
	}
	seqLine := p.scanner.Text()

	if !p.scanner.Scan() {
		return nil, fmt.Errorf("unexpected EOF after sequence line")
	}
	sepLine := p.scanner.Text()

	if !p.scanner.Scan() {
		return nil, fmt.Errorf("unexpected EOF after separator line")
	}
	qualLine := p.scanner.Text()

	// Validation
	if idLine[0] != '@' {
		return nil, fmt.Errorf("expected id line to start with '@', got '%s'", idLine)
	}
	if sepLine[0] != '+' {
		return nil, fmt.Errorf("expected separator line to start with '+', got '%s'", sepLine)
	}
	if len(seqLine) != len(qualLine) {
		return nil, fmt.Errorf("sequence and quality length mismatch (%d vs %d)", len(seqLine), len(qualLine))
	}

	// Populate and return the record.
	record := &FastqRecord{
		ID:   idLine[1:], // Remove leading '@'
		Seq:  []byte(seqLine),
		Qual: []byte(qualLine),
	}

	return record, nil

}
