package ui

import "github.com/charmbracelet/lipgloss"

// renderGraph renders the sparkline graph for the selected interface.
// Full implementation in Task 8.
func renderGraph(_ Model, w, h int) string {
	return lipgloss.NewStyle().Width(w).Height(h).Render("Graph panel")
}
