package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/guillechuma/bio-tui/internal/fastq"
)

func main() {
	// 1. Open the test FASTQ file.
	f, err := os.Open("testdata/reads.fastq")
	if err != nil {
		log.Fatalf("could not open test file: %v", err)
	}
	defer f.Close()

	// 2. Create a new FASTQ parser.
	parser := fastq.NewParser(f)
	fmt.Println("--- Reading FASTQ records ---")

	// 3. Loop through all the records in the file.
	for {
		record, err := parser.Next()
		if err != nil {
			if err == io.EOF {
				break // We've reached the end of the file.
			}
			log.Fatalf("parser failed: %v", err)
		}

		// 4. If we got a valid record, print its details.
		fmt.Printf("Successfully parsed record!\n")
		fmt.Printf("\tID:   %s\n", record.ID)
		fmt.Printf("\tSeq:  %s (len: %d)\n", string(record.Seq), len(record.Seq))
		fmt.Printf("\tQual: %s (len: %d)\n", string(record.Qual), len(record.Qual))
		fmt.Printf("\tMean Quality: %.2f\n", record.MeanQual())
		fmt.Printf("\tGC Content:   %.2f%%\n", record.GCContent()*100)
	}

	fmt.Println("--- Done ---")
}
