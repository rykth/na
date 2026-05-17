package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderIfList renders the left-side scrollable interface list panel.
func renderIfList(m Model, height int) string {
	names := m.store.Names()

	// compute scroll window so selected is always visible
	offset := 0
	if m.selected >= height {
		offset = m.selected - height + 1
	}

	var sb strings.Builder
	for i := offset; i < len(names) && i < offset+height; i++ {
		sb.WriteString(renderIfLine(m, names[i], i == m.selected))
		sb.WriteByte('\n')
	}

	content := strings.TrimRight(sb.String(), "\n")
	return lipgloss.NewStyle().Width(listWidth).Height(height).Render(content)
}

// renderIfLine renders a single interface row.
func renderIfLine(m Model, name string, selected bool) string {
	iface := m.store.Get(name)

	prefix := "  "
	if selected {
		prefix = "> "
	}

	state := "[??]"
	var rowStyle lipgloss.Style

	if iface != nil {
		switch strings.ToLower(iface.OperState) {
		case "up":
			if iface.Stats.RxBytes.Rate > 0 || iface.Stats.TxBytes.Rate > 0 {
				rowStyle = styleUp.Bold(true)
			} else {
				rowStyle = styleDim
			}
			state = "[UP]"
		case "down":
			rowStyle = styleDown
			state = "[DOWN]"
		default:
			rowStyle = styleDim
			state = "[" + strings.ToUpper(iface.OperState) + "]"
		}
	}

	rx, tx := 0.0, 0.0
	if iface != nil {
		rx = iface.Stats.RxBytes.Rate
		tx = iface.Stats.TxBytes.Rate
	}

	rxStr := formatRate(rx, m.useBits, m.useSI)
	txStr := formatRate(tx, m.useBits, m.useSI)

	// truncate name to keep line within listWidth
	displayName := name
	if len(displayName) > 8 {
		displayName = displayName[:8]
	}

	line := fmt.Sprintf("%s%-8s ↓ %-10s ↑ %-10s %s",
		prefix, displayName, rxStr, txStr, state)

	if selected {
		return styleSelected.Render(line)
	}
	return rowStyle.Render(line)
}
