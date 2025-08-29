package fastq

// FastqRecord holds the data for a single FASTQ entry.
type FastqRecord struct {
	ID   string // Sequence ID
	Seq  []byte // The raw sequence data
	Qual []byte // The quality scores for the sequence
}

// MeanQual calculates the average Phred quality score for the read
func (r *FastqRecord) MeanQual() float64 {
	if len(r.Qual) == 0 {
		return 0.0
	}

	var totalQuality int
	for _, q := range r.Qual {
		// Convert ASCII char to Phred score (Phred+33 encoding)
		score := int(q) - 33
		totalQuality += score
	}

	return float64(totalQuality) / float64(len(r.Qual))
}

// GCContent calculates the percentage of Guanine (G) and Cytosine (C)
// bases in the sequence.
func (r *FastqRecord) GCContent() float64 {
	if len(r.Seq) == 0 {
		return 0.0
	}

	gcCount := 0
	for _, base := range r.Seq {
		// Make the check case-insensitive.
		switch base {
		case 'G', 'g', 'C', 'c':
			gcCount++
		}
	}

	return float64(gcCount) / float64(len(r.Seq))
}
