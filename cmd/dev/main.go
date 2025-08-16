package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/guillechuma/bio-tui/internal/adapters/fasta"
)

func main() {
	// Open the hardcoded test file.
	filePath := "testdata/test.fasta"
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("could not open test file: %v", err)
	}
	defer f.Close() // Make sure the file gets closed when main() exits.

	var fileReader io.Reader = f

	// Conditionally wrap the reader if the file is gzipped
	if strings.HasSuffix(filePath, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			log.Fatalf("could not create gzip reader: %v", err)
		}
		defer gz.Close()
		fileReader = gz
	}

	// It assumes you have a NewParser function and a Next method.
	parser := fasta.NewParser(fileReader)

	start := time.Now() // <-- 1. Record start time

	// Loop through all the records in the file.
	for {
		record, err := parser.Next()
		if err != nil {
			// io.EOF is the error the parser will return when it's done.
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("parser failed: %v", err)
		}

		// If we got here, we have a valid record!
		fmt.Printf("Successfully parsed record!\n")
		fmt.Printf("\tID:  %s\n", record.ID)
		fmt.Printf("\tSeq: %s\n", string(record.Seq))
		fmt.Printf("\tGC Content: %.2f%%\n", record.GCContent()*100)
	}

	duration := time.Since(start) // <-- 2. Calculate elapsed time

	fmt.Println("Done.")
	fmt.Printf("Parsing took: %s\n", duration) // <-- 3. Print it!
}
