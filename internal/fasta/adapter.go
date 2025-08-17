package fasta

import (
	"fmt"

	"github.com/guillechuma/bio-tui/internal/adapter"
)

// FastaAdapter satisfies the adapter.Reader interface for FASTA files.
// It uses an IndexedFastaReader to provide fast, random access.
type FastaAdapter struct {
	reader *IndexedReader
}

// Open initializes the reader by loading the FASTA and its index.
func (a *FastaAdapter) Open(spec adapter.OpenSpec) error {
	r, err := NewIndexedReader(spec.Path) // Creates the local IndexedReader
	if err != nil {
		return err
	}
	a.reader = r
	return nil
}

// Close releases the underlying file handle.
func (a *FastaAdapter) Close() error {
	if a.reader != nil {
		return a.reader.Close()
	}
	return nil
}

// Capabilities reports that this adapter can look up symbols and read regions.
func (a *FastaAdapter) Capabilities() adapter.Capability {
	return adapter.CapSymbols | adapter.CapRegions
}

// ListSymbols returns a slice of all sequence IDs and their lengths.
func (a *FastaAdapter) ListSymbols() ([]adapter.Symbol, error) {
	symbols := make([]adapter.Symbol, 0, len(a.reader.Index)) // Array length in memory record length of index
	for name, record := range a.reader.Index {
		symbols = append(symbols, adapter.Symbol{
			Name:   name,
			Length: record.Length,
		})
	}
	return symbols, nil
}

// LookupSymbol finds a sequence by its ID and returns its full region.
func (a *FastaAdapter) LookupSymbol(sym string) (adapter.Region, error) {
	// We need the length of the sequence, which is in the index.
	// We'll expose the index map for this.
	indexRecord, ok := a.reader.Index[sym]
	if !ok {
		return adapter.Region{}, fmt.Errorf("symbol '%s' not found in FASTA index", sym)
	}

	reg := adapter.Region{
		Ref:   sym,
		Start: 1,
		End:   indexRecord.Length,
	}

	return reg, nil
}

// Region fetches the sequence data for a specific genomic region.
func (a *FastaAdapter) Region(reg adapter.Region) (adapter.Slice, error) {
	// First, we fetch the entire sequence record for the given reference (e.g., 'chr1').
	record, err := a.reader.Fetch(reg.Ref)
	if err != nil {
		return adapter.Slice{}, err
	}

	// Then, we use the Slice method we already built on FastaRecord
	// to extract the specific subsequence.
	subsequence, err := record.Slice(int(reg.Start), int(reg.End))
	if err != nil {
		return adapter.Slice{}, err
	}

	slice := adapter.Slice{
		Sequence: subsequence,
	}

	return slice, nil
}

// IterRows is not applicable to FASTA files in a meaningful way,
// so we return an "unsupported" error.
func (a *FastaAdapter) IterRows(ch chan<- []string, stop <-chan struct{}) error {
	return fmt.Errorf("IterRows is not supported by the FastaAdapter")
}
