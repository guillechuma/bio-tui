package index

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FaiRecord holds the index information for a single sequence in a FASTA file.
type FaiRecord struct {
	Name      string // Name of the sequence
	Length    int64  // Total length of the sequence in bases
	Offset    int64  // Byte offset in the file where the sequence starts
	LineBases int64  // Number of bases per line
	LineBytes int64  // Number of bytes per line (including newline)
}

// ParseFai reads a .fai file and returns a map of sequence names to their index records.
func ParseFai(path string) (map[string]FaiRecord, error) {
	// Open index file
	indexFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open index file (.fai): %w", err)
	}
	defer indexFile.Close()

	// Read the index line-by-line and store it in a map.
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

	return index, scanner.Err()
}
