package ui

import "github.com/charmbracelet/lipgloss"

// renderIfList renders the left-side interface list panel.
// Full implementation in Task 6.
func renderIfList(_ Model, h int) string {
	return lipgloss.NewStyle().Width(listWidth).Height(h).Render("Interface list")
}
