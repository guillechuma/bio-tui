package main

import (
	"fmt"
	"log"
	"time"

	"github.com/guillechuma/bio-tui/internal/adapter"
	"github.com/guillechuma/bio-tui/internal/fasta"
)

func main() {
	// Open the hardcoded test file.
	filePath := "testdata/test.fasta"
	reader, err := fasta.NewIndexedReader(filePath)
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

	// 1. Create an instance of your concrete adapter.
	fa := &fasta.FastaAdapter{}

	// 2. Assign it to a variable of the INTERFACE type.
	//    This proves your adapter has the right "shape".
	var rdr adapter.Reader = fa

	// 3. Open the file using the interface method.
	spec := adapter.OpenSpec{Path: "testdata/test.fasta"}
	err = rdr.Open(spec)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// 4. Use the interface methods to get data.
	fmt.Println("Looking up 'header2'...")
	region, err := rdr.LookupSymbol("header2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found 'header2' at region: %s:%d-%d\n", region.Ref, region.Start, region.End)

	fmt.Println("\nFetching a sub-region from 'header2'...")
	subRegion := adapter.Region{Ref: "header2", Start: 3, End: 10}
	slice, err := rdr.Region(subRegion)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sequence for %s:%d-%d is: %s\n", subRegion.Ref, subRegion.Start, subRegion.End, string(slice.Sequence))
}
