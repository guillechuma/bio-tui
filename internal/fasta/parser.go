package fasta

import (
	"bufio"
	"bytes"
	"fmt"
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
	record.Type = InferSequenceType(record.Seq)

	// Now, validate the record based on its inferred type.
	if !record.Validate() {
		return nil, fmt.Errorf("record %s contains invalid characters for inferred type", record.ID)
	}

	return record, nil
}

// InferSequenceType analyzes a byte slice and returns the likely sequence type.
func InferSequenceType(seq []byte) SequenceType {
	// Flags to track which character sets we've seen.
	hasT := false
	hasU := false
	hasProteinChars := false

	for _, base := range seq {
		switch base {
		case 'T', 't':
			hasT = true
		case 'U', 'u':
			hasU = true
			// Check for amino acids that are not also ambiguity codes for DNA/RNA
		case 'E', 'e', 'F', 'f', 'I', 'i', 'L', 'l', 'P', 'p', 'Q', 'q', 'Z', 'z', 'X', 'x', '*':
			hasProteinChars = true
		}
	}

	if hasProteinChars {
		return Protein
	}
	if hasT && hasU {
		return UnknownSequence
	}
	if hasU {
		return RNA
	}
	// If it has a T, or if it has neither T nor U (e.g. "ACGN"),
	// we default to DNA, as it's the most common type.
	return DNA
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

// isValidDNA checks if a sequence contains only valid DNA characters,
// including ambiguity codes.
func isValidDNA(seq []byte) bool {
	for _, base := range seq {
		switch base {
		// Standard Bases
		case 'A', 'a', 'C', 'c', 'G', 'g', 'T', 't':
		// Uracil is sometimes found
		case 'U', 'u':
		// Ambiguity Codes
		case 'R', 'r', // A or G
			'Y', 'y', // C or T
			'S', 's', // G or C
			'W', 'w', // A or T
			'K', 'k', // G or T
			'M', 'm', // A or C
			'B', 'b', // C, G, or T
			'D', 'd', // A, G, or T
			'H', 'h', // A, C, or T
			'V', 'v', // A, C, or G
			'N', 'n': // Any
		// Gaps are also common
		case '-':
		default:
			return false // Invalid character found
		}
	}
	return true
}

// isValidRNA checks if a sequence contains only valid RNA characters,
// including ambiguity codes.
func isValidRNA(seq []byte) bool {
	for _, base := range seq {
		switch base {
		// Standard Bases
		case 'A', 'a', 'C', 'c', 'G', 'g', 'U', 'u':
		// Ambiguity Codes (same as DNA)
		case 'R', 'r', 'Y', 'y', 'S', 's', 'W', 'w',
			'K', 'k', 'M', 'm', 'B', 'b', 'D', 'd',
			'H', 'h', 'V', 'v', 'N', 'n':
		// Gaps
		case '-':
		default:
			return false // Invalid character found
		}
	}
	return true
}

// isValidProtein checks if a sequence contains only valid amino acid characters.
func isValidProtein(seq []byte) bool {
	for _, aa := range seq {
		switch aa {
		// Standard 20 Amino Acids
		case 'A', 'a', 'C', 'c', 'D', 'd', 'E', 'e', 'F', 'f',
			'G', 'g', 'H', 'h', 'I', 'i', 'K', 'k', 'L', 'l',
			'M', 'm', 'N', 'n', 'P', 'p', 'Q', 'q', 'R', 'r',
			'S', 's', 'T', 't', 'V', 'v', 'W', 'w', 'Y', 'y':
		// Ambiguity and Special Characters
		case 'B', 'b', // Aspartic acid or Asparagine
			'Z', 'z', // Glutamic acid or Glutamine
			'X', 'x', // Any amino acid
			'*', // Stop codon
			'-': // Gap
		default:
			return false // Invalid character found
		}
	}
	return true
}
