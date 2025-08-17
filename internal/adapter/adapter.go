package adapter

// OpenSpec provides all the necessary information to open an adapter.
type OpenSpec struct {
	Path  string            // The primary file path (e.g., sample.bam, genes.gff)
	Index string            // An optional, explicit path to an index file (e.g., .tbi, .bai)
	Aux   map[string]string // A map for any other auxiliary files the adapter might need.
}

// Region defines a genomic interval. It is half-open: [Start, End).
type Region struct {
	Ref   string
	Start int64
	End   int64
}

// Capability is a bitmask to report what an adapter can do.
type Capability uint32

// Set of possible capabilities
const (
	CapRegions  Capability = 1 << iota // Supports region-limited queries.
	CapIterRows                        // Can stream all records (for tables).
	CapCoverage                        // Can compute coverage track
	CapPileup                          // Can compute a pileup track.
	CapSymbols                         // Can look up features by name (e.g., gene ID).
)

// Slice contains all the data for a given genomic region.
// This is the primary data structure returned to the UI for rendering.
type Slice struct {
	// Hold sequence data for the region.
	Sequence []byte
}

// Reader is the universal interface for all file type adapters.
// It defines a standard contract for the TUI to interact with data sources.
type Reader interface {
	Open(spec OpenSpec) error
	Close() error
	Capabilities() Capability
	Region(reg Region) (Slice, error)
	LookupSymbol(sym string) (Region, error)
	IterRows(ch chan<- []string, stop <-chan struct{}) error
}
