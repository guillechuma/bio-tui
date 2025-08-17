package fasta

import (
	"bytes"
	"testing"
)

func TestFastaRecord_Slice(t *testing.T) {
	// Set up test case
	record := FastaRecord{
		ID:  "test_seq",
		Seq: []byte("ACGTACGTAC"), // Length is 10
	}

	// Inputs and expected outputs
	start, end := 3, 7 // 1-based coordinates
	expected := []byte("GTACG")

	// Run slice
	actual, err := record.Slice(start, end)
	if err != nil {
		t.Fatalf("Slice() returned an unexpected error: %v", err)
	}

	// Assert result
	if !bytes.Equal(actual, expected) {
		t.Errorf("Slice() failed: expected %s, got %s", expected, actual)
	}

}
