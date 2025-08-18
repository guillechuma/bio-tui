package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guillechuma/bio-tui/internal/adapter"
	"github.com/guillechuma/bio-tui/internal/fasta"
	"github.com/guillechuma/bio-tui/internal/ui"
)

func main() {
	// 1. Check for a command-line argument for the file path.
	if len(os.Args) < 2 {
		fmt.Println("Usage: bio-tui <fasta-file>")
		os.Exit(1)
	}
	filePath := os.Args[1]

	// 2. Set up your FastaAdapter.
	fa := &fasta.FastaAdapter{}
	var reader adapter.Reader = fa

	// 3. Open the file and get the list of symbols.
	spec := adapter.OpenSpec{Path: filePath}
	if err := reader.Open(spec); err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer reader.Close()

	symbols, err := reader.ListSymbols()
	if err != nil {
		log.Fatalf("Error listing symbols: %v", err)
	}

	// 4. Create the TUI model with the data.
	model := ui.NewModel(symbols, reader)

	// 5. Create and run the Bubble Tea program.
	// Using WithAltScreen restores the terminal to its original state on exit.
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
