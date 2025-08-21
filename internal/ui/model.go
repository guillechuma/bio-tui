// This file defines the state and behavior of the TUI application

package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guillechuma/bio-tui/internal/adapter"
)

type focusState int

const (
	focusList focusState = iota
	focusViewport
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
	adapter      adapter.Reader // Store the adapter to fetch data
	list         list.Model
	viewport     viewport.Model // For the sequence viewer
	styles       Styles
	focus        focusState
	currentSlice adapter.Slice
	quitting     bool
}

// NewModel creates and returns a new TUI model, initialized with the sequence symbols.
func NewModel(symbols []adapter.Symbol, reader adapter.Reader) Model {
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

	// The viewport is not visible yet, but we initialize it.
	vp := viewport.New(0, 0)

	return Model{
		adapter:  reader,
		list:     ls,
		viewport: vp,
		styles:   NewStyles(), // Initialize styles
		focus:    focusList,   // <-- Start with the list focused
	}
}

// Init is the first command that's run when the program starts.
func (m Model) Init() tea.Cmd {
	return nil // No initial command needed.
}

// Update is the main event loop. It handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// Handle window resize events.
	case tea.WindowSizeMsg:
		// This is the total screen width and height.
		screenWidth := msg.Width
		screenHeight := msg.Height

		// --- Calculate Pane Widths ---
		// Get the horizontal "overhead" (borders + padding) from the styles.
		listStyle := m.styles.Active
		viewportStyle := m.styles.Inactive
		listOverhead := listStyle.GetHorizontalFrameSize()
		viewportOverhead := viewportStyle.GetHorizontalFrameSize()

		// Define the width for the list pane (e.g., 1/3 of the screen).
		listPaneWidth := screenWidth / 3
		// The viewport gets the rest of the space.
		viewportPaneWidth := screenWidth - listPaneWidth

		// --- Set Component Sizes ---
		// The list's content area is its pane width minus its style's overhead.
		m.list.SetSize(listPaneWidth-listOverhead, screenHeight-2)

		// The viewport's content area is its pane width minus its style's overhead.
		m.viewport.Width = viewportPaneWidth - viewportOverhead
		m.viewport.Height = screenHeight - 2

		// Re-wrap the content from the stored slice.
		// The IsEmpty check for the slice is more robust than checking for an empty string.
		if m.currentSlice.Sequence != nil {
			wrappedSequence := m.wrapSequence(string(m.currentSlice.Sequence), 1)
			m.viewport.SetContent(wrappedSequence)
		}

	// Handle key presses.
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "tab":
			// Switch focus between the list and the viewport.
			if m.focus == focusList {
				m.focus = focusViewport
			} else {
				m.focus = focusList
			}
			return m, nil
		}
	}

	// --- Component-Specific Message Routing ---
	// Instead of sending messages to both components, we now route them
	// to only the component that has focus.
	switch m.focus {
	case focusList:
		// The list is focused.
		beforeIndex := m.list.Index()
		m.list, cmd = m.list.Update(msg)
		if m.list.Index() != beforeIndex {
			// The selection changed, so update the viewport content.
			return m, m.updateViewportContent()
		}
	case focusViewport:
		// The viewport is focused, so only it should receive updates.
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

// view renders de TUI.
func (m Model) View() string {
	if m.quitting {
		return "Bye!\n"
	}
	// --- Dynamic Style Assignment ---
	// Declare two local style variables.
	var listStyle, viewportStyle lipgloss.Style

	// Check which pane has focus and assign the Active/Inactive
	// styles accordingly.
	if m.focus == focusList {
		listStyle = m.styles.Active
		viewportStyle = m.styles.Inactive
	} else { // focus == focusViewport
		listStyle = m.styles.Inactive
		viewportStyle = m.styles.Active
	}
	// Render the list and viewport into their own strings.
	listView := listStyle.Render(m.list.View())
	viewportView := viewportStyle.Render(m.viewport.View())
	statsView := m.renderStatsPanel() // A new helper function

	// stack the viewport and stats panel vertically.
	rightPane := lipgloss.JoinVertical(lipgloss.Left, viewportView, statsView)

	// Use lipgloss to join them horizontally.
	return lipgloss.JoinHorizontal(lipgloss.Top, listView, rightPane)
}

// wrapSequence wraps a DNA sequence to fit within the viewport width,
// prepending each line with its genomic coordinate.
func (m *Model) wrapSequence(sequence string, startCoord int) string {
	if m.viewport.Width <= 0 {
		return sequence
	}

	// 1. Get the style that will be used for the viewport pane.
	style := m.styles.Inactive

	// 2. Ask the style for its horizontal padding.
	padding := style.GetHorizontalPadding()
	// Define the width of our position indicator margin (e.g., "1234567890 ").
	marginWidth := 11

	// The space available for the sequence is the viewport's inner width minus our margin.
	availableWidth := m.viewport.Width - padding
	lineWidth := availableWidth - marginWidth

	if lineWidth <= 0 {
		return sequence
	}

	var wrapped strings.Builder
	currentCoord := startCoord
	for i := 0; i < len(sequence); i += lineWidth {
		end := min(i+lineWidth, len(sequence))

		if i > 0 {
			wrapped.WriteString("\n")
		}

		// Prepend the formatted coordinate.
		// The `%-10d` format right-pads the number with spaces to a width of 10.
		coordMargin := fmt.Sprintf("%-10d", currentCoord)
		wrapped.WriteString(coordMargin)

		wrapped.WriteString(sequence[i:end])

		// Increment our coordinate for the next line.
		currentCoord += lineWidth
	}

	return wrapped.String()
}

// updateViewportContent is a new helper function to fetch and set the viewport data.
func (m *Model) updateViewportContent() tea.Cmd {
	// Get the currently selected item.
	selectedItem, ok := m.list.SelectedItem().(item)
	if !ok {
		return nil
	}

	// Use the adapter to fetch the full sequence.
	region := adapter.Region{Ref: selectedItem.symbol.Name, Start: 1, End: selectedItem.symbol.Length}
	slice, err := m.adapter.Region(region)
	if err != nil {
		m.viewport.SetContent(fmt.Sprintf("Error: %v", err))
		return nil
	}

	// 1. Store the entire generic Slice object in model
	m.currentSlice = slice

	// 2. Wrap and set the content from the slice's Sequence field.
	wrappedSequence := m.wrapSequence(string(m.currentSlice.Sequence), 1)
	m.viewport.SetContent(wrappedSequence)

	// Go back to the top of the viewport every time the content changes.
	m.viewport.GotoTop()

	return nil
}

func (m Model) renderStatsPanel() string {
	style := m.styles.Inactive

	stats := m.currentSlice.Stats
	if len(stats) == 0 {
		return style.Render("")
	}

	// Get all the keys from the map.
	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}

	// Sort the keys alphabetically for a consistent order.
	sort.Strings(keys)

	// Build the content string by iterating over the sorted keys.
	var contentBuilder strings.Builder
	for _, key := range keys {
		value := stats[key]
		// Left-align the key, right-align the value.
		line := lipgloss.JoinHorizontal(lipgloss.Left,
			fmt.Sprintf("%-12s", key), // Pad the key for alignment
			lipgloss.NewStyle().Width(m.viewport.Width-12).Align(lipgloss.Right).Render(value),
		)
		contentBuilder.WriteString(line)
		contentBuilder.WriteString("\n")
	}

	return style.Render(strings.TrimSpace(contentBuilder.String()))
}
