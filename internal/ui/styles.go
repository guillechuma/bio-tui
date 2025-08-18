package ui

import "github.com/charmbracelet/lipgloss"

// Struct to hold all styles
type Styles struct {
	Base,
	Active,
	Inactive lipgloss.Style
}

// NewStyles creates a new Styles struct with default settings.
func NewStyles() Styles {
	s := Styles{}
	// A base style for both panes
	s.Base = lipgloss.NewStyle().
		Padding(1, 2).
		BorderForeground(lipgloss.Color("240")) // A nice dim gray

	// Style for the active pane
	s.Active = s.Base.Copy().
		Border(lipgloss.DoubleBorder(), true).
		BorderForeground(lipgloss.Color("205")) // A vibrant pink/magenta

	// Style for the inactive pane
	s.Inactive = s.Base.Copy().
		Border(lipgloss.NormalBorder(), true)

	return s
}
