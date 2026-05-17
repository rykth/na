package ui

import "github.com/charmbracelet/lipgloss"

// renderStats renders the right-side stats panel for the selected interface.
// Full implementation in Task 7.
func renderStats(_ Model, w, h int) string {
	return lipgloss.NewStyle().Width(w).Height(h).Render("Stats panel")
}
