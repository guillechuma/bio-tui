package index

import (
	"bufio"
	"fmt"
	"io"
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

// BuildFai creates a .fai index for a given FASTA file.
// validating for consistent line lengths.
func BuildFai(fastaPath string) error {
	// Open the input fasta
	inFile, err := os.Open(fastaPath)
	if err != nil {
		return fmt.Errorf("could not open fasta file to build index: %w", err)
	}
	defer inFile.Close()

	// Create the output .fai file.
	outFile, err := os.Create(fastaPath + ".fai")
	if err != nil {
		return fmt.Errorf("could not create .fai file: %w", err)
	}
	defer outFile.Close()

	// use a bufio.Reader for line-by-line reading with offset tracking.
	reader := bufio.NewReader(inFile)
	var byteOffset int64 = 0

	// State variables for the current sequence record
	var current struct {
		name                   string
		length                 int64
		offset                 int64
		lineBases              int64
		lineBytes              int64
		isShortLineEncountered bool // <-- New flag to track state
	}

	for {
		// Read one line, keeping track of how many bytes were consumed.
		line, err := reader.ReadString('\n')
		lineBytesRead := int64(len(line))
		// Trim whitespace from the line for processing.
		trimmedLine := strings.TrimSpace(line)

		if err != nil && err != io.EOF {
			return err // Handle unexpected errors
		}

		// Check if it's a header line.
		if strings.HasPrefix(trimmedLine, ">") {
			// If we were already tracking a sequence, we must write its entry first.
			if current.name != "" {
				fmt.Fprintf(outFile, "%s\t%d\t%d\t%d\t%d\n",
					current.name, current.length, current.offset, current.lineBases, current.lineBytes)
			}
			// Start a new record.
			current.name = strings.TrimSpace(trimmedLine[1:])
			current.length = 0
			current.lineBases = 0
			current.offset = byteOffset + lineBytesRead // The first base is on the *next* line
			current.isShortLineEncountered = false      // Reset the flag.
		} else if current.name != "" && len(trimmedLine) > 0 {
			// This is a sequence line.
			// **VALIDATION 1: Check if we've already seen a short line.**
			if current.isShortLineEncountered {
				return fmt.Errorf("format error in sequence '%s': unexpected sequence data after a short line", current.name)
			}

			if current.lineBases == 0 {
				// First sequence line; this sets the standard.
				current.lineBases = int64(len(trimmedLine))
				current.lineBytes = lineBytesRead
			} else {
				// **VALIDATION 2: Check for inconsistent line length.**
				if int64(len(trimmedLine)) != current.lineBases {
					// This line is shorter than the standard. It must be the last one.
					if int64(len(trimmedLine)) > current.lineBases {
						return fmt.Errorf("different line length in sequence '%s'", current.name)
					}
					// It's a short line. Set the flag.
					current.isShortLineEncountered = true
				}
			}
			// This is a sequence line for the current record.
			current.length += int64(len(trimmedLine))
		}

		byteOffset += lineBytesRead
		if err == io.EOF {
			break
		}
	}
	// Write the very last record after the loop finishes.
	if current.name != "" {
		fmt.Fprintf(outFile, "%s\t%d\t%d\t%d\t%d\n",
			current.name, current.length, current.offset, current.lineBases, current.lineBytes)
	}

	return nil
}
