package fasta

import "fmt"

type SequenceType int

const (
	UnknownSequence SequenceType = iota // Default 0
	DNA
	RNA
	Protein
)

type FastaRecord struct {
	ID          string
	Description string
	Seq         []byte
	Type        SequenceType
}

// Validate checks a record's sequence against the alphabet of its inferred type.
func (r *FastaRecord) Validate() bool {
	switch r.Type {
	case DNA:
		return isValidDNA(r.Seq)
	case RNA:
		return isValidRNA(r.Seq)
	case Protein:
		return isValidProtein(r.Seq)
	default:
		return false
	}
}

func (r *FastaRecord) Slice(start, end int) ([]byte, error) {
	// 1. Validate coordinates
	seqLen := len(r.Seq)
	if start < 1 || end > seqLen || start > end {
		return nil, fmt.Errorf("invalid slice coordinates: start %d, end %d for sequence of length %d", start, end, seqLen)
	}

	// 2. Convert 1-based (human) to 0-based (Go) ---
	// The 0-based start is simply start - 1.
	// The 0-based end for Go's half-open slicing is just `end`.
	zeroBasedStart := start - 1
	zeroBasedEnd := end

	// 3. Slice
	subsequence := r.Seq[zeroBasedStart:zeroBasedEnd]

	return subsequence, nil
}

// GCContent calculates the percentage of Guanine (G) and Cytosine (C)
// bases in the sequence.
func (r *FastaRecord) GCContent() float64 {
	// 1. Handle edge case of an empty sequence to avoid div by 0
	if len(r.Seq) == 0 {
		return 0.0
	}

	// 2. Initialize a counter for G and C bases
	gcCount := 0

	// 3. Loop through the sequence
	for _, base := range r.Seq {
		// Case-insentive.
		switch base {
		case 'G', 'g', 'C', 'c':
			gcCount++
		}
	}

	// 4. GC calculation.
	return float64(gcCount) / float64(len(r.Seq))
}
