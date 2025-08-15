package main

import (
	"fmt"
	"log"
	"os"

	"github.com/guillechuma/bio-tui/internal/adapters/fasta"
)

func main() {
	// Open the hardcoded test file.
	f, err := os.Open("testdata/test.fasta")
	if err != nil {
		log.Fatalf("could not open test file: %v", err)
	}
	defer f.Close() // Make sure the file gets closed when main() exits.

	// It assumes you have a NewParser function and a Next method.
	parser := fasta.NewParser(f)

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
	}

	fmt.Println("Done.")
}
