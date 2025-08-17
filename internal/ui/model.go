// This file defines the state and behavior of the TUI application

package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/guillechuma/bio-tui/internal/adapter"
)

// item represents a single sequence in our list. It needs to satisfy
// the list.Item interface for the bubbles/list component.
type item struct {
	symbol adapter.Symbol
}

// Title is what's shown in the list item's main line
func (i item) Title() string { return i.symbol.Name }

// Description is shown below the title.
func (i item) Description() string { return fmt.Sprintf("%d bp", i.symbol.Length) }

// FilterValue is the string the list will filter against.
func (i item) FilterValue() string { return i.symbol.Name }

// Model holds the state of our TUI application.
type Model struct {
	list     list.Model
	quitting bool
}

// NewModel creates and returns a new TUI model, initialized with the sequence symbols.
func NewModel(symbols []adapter.Symbol) Model {
	// 1. Convert our []adapter.Symbol into a []list.Item for the component.
	items := make([]list.Item, len(symbols))
	for i, sym := range symbols {
		items[i] = item{symbol: sym}
	}

	// 2. Setup the list component.
	ls := list.New(items, list.NewDefaultDelegate(), 0, 0)
	ls.Title = "Fasta Sequences"
	ls.SetShowStatusBar(true)
	ls.SetFilteringEnabled(true)

	return Model{list: ls}
}

// Init is the first command that's run when the program starts.
func (m Model) Init() tea.Cmd {
	return nil // No initial command needed.
}

// Update is the main event loop. It handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle window resize events.
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	// Handle key presses.
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Pass all other messages to the list's own Update function.
	// This handles scrolling, filtering, etc.
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// view renders de TUI.
func (m Model) View() string {
	if m.quitting {
		return "Bye!\n"
	}
	// Just render the list component.
	return m.list.View()
}
