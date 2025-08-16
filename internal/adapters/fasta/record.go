package fasta

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
		return isValidProtein(R.seq)
	default:
		return false
	}
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

	// 4. Perform the final calculation.
	// We must convert the numbers to float64 *before* dividing
	// to ensure we get a floating-point result
	return float64(gcCount) / float64(len(r.Seq))
}
