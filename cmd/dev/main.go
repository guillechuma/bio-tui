package main

import (
	"fmt"
	"log"
	"time"

	"github.com/guillechuma/bio-tui/internal/index"
)

func main() {
	// Open the hardcoded test file.
	filePath := "testdata/test.fasta"
	reader, err := index.NewIndexedFastaReader(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close() // Make sure the file gets closed when main() exits.

	start := time.Now() // <-- 1. Record start time
	fmt.Println("Fetching record 'header2'...")
	record, err := reader.Fetch("header2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully fetched record!\n")
	fmt.Printf("\tID:  %s\n", record.ID)
	fmt.Printf("\tSeq: %s\n", string(record.Seq))
	fmt.Printf("\tGC Content: %.2f%%\n", record.GCContent()*100)

	// // Loop through all the records in the file.
	// for {
	// 	record, err := parser.Next()
	// 	if err != nil {
	// 		// io.EOF is the error the parser will return when it's done.
	// 		if err.Error() == "EOF" {
	// 			break
	// 		}
	// 		log.Fatalf("parser failed: %v", err)
	// 	}

	// 	// If we got here, we have a valid record!
	// 	fmt.Printf("Successfully parsed record!\n")
	// 	fmt.Printf("\tID:  %s\n", record.ID)
	// 	fmt.Printf("\tSeq: %s\n", string(record.Seq))
	// 	fmt.Printf("\tGC Content: %.2f%%\n", record.GCContent()*100)
	// }

	duration := time.Since(start) // <-- 2. Calculate elapsed time

	fmt.Println("Done.")
	fmt.Printf("Fetching took: %s\n", duration) // <-- 3. Print it!
}
