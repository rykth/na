package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var styleHelpBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1, 3)

// renderHelp renders the full-screen help overlay.
func renderHelp(m Model) string {
	bindings := []struct{ key, desc string }{
		{"↑ / k", "Select previous interface"},
		{"↓ / j", "Select next interface"},
		{"g", "Cycle graph: off → RX → RX+TX"},
		{"d", "Toggle detailed counters"},
		{"b", "Toggle bytes / bits display"},
		{"s", "Toggle SI / IEC units"},
		{"l", "Toggle interface list"},
		{"r", "Reset rate accumulators"},
		{"+ / -", "Adjust poll interval (±0.5 s)"},
		{"h / ?", "Toggle this help"},
		{"q / Ctrl-C", "Quit"},
	}

	var sb strings.Builder
	sb.WriteString("na — Network Analyzer\n")
	sb.WriteString(strings.Repeat("─", 44) + "\n\n")
	for _, b := range bindings {
		sb.WriteString(fmt.Sprintf("  %-14s  %s\n", b.key, b.desc))
	}
	sb.WriteString("\n" + strings.Repeat("─", 44) + "\n")
	sb.WriteString("Press any key to dismiss")

	box := styleHelpBox.Render(sb.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
