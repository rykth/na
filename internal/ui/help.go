package ui

import "github.com/charmbracelet/lipgloss"

// renderHelp renders the full-screen help overlay.
// Full implementation in Task 9.
func renderHelp(m Model) string {
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		"Help overlay — press any key to dismiss")
}
