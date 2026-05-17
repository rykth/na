package ui

import "github.com/charmbracelet/lipgloss"

// renderDetails renders the expanded per-counter details panel.
// Full implementation in Task 7.
func renderDetails(_ Model, w, h int) string {
	return lipgloss.NewStyle().Width(w).Height(h).Render("Details panel")
}
