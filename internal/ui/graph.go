package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	graphHeight = 6 // data rows per graph
	graphLines  = graphHeight + 2 // +1 title +1 bottom axis
	yAxisW      = 1 // width of the '│'/'└' column
)

var blocks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// renderGraph renders RX (and optionally TX) sparkline graphs.
// showGraph: 0=off, 1=RX only, 2=RX+TX.
func renderGraph(m Model, w, h int) string {
	if m.showGraph == 0 {
		return ""
	}

	names := m.store.Names()
	if len(names) == 0 || m.selected >= len(names) {
		return lipgloss.NewStyle().Width(w).Height(h).Render("")
	}

	iface := m.store.Get(names[m.selected])
	if iface == nil || iface.History == nil {
		return lipgloss.NewStyle().Width(w).Height(h).Render("")
	}

	gw := w - yAxisW
	if gw < 1 {
		gw = 1
	}

	rx, tx := iface.History.Samples(gw)

	var sb strings.Builder
	sb.WriteString(drawGraph("RX", rx, gw, m.useBits, m.useSI))
	if m.showGraph == 2 {
		sb.WriteByte('\n')
		sb.WriteString(drawGraph("TX", tx, gw, m.useBits, m.useSI))
	}

	return lipgloss.NewStyle().Width(w).Height(h).Render(sb.String())
}

// drawGraph renders a single labelled sparkline graph for the given samples.
func drawGraph(label string, samples []float64, gw int, useBits, useSI bool) string {
	peak := 0.0
	for _, s := range samples {
		if s > peak {
			peak = s
		}
	}
	divisor := peak
	if divisor == 0 {
		divisor = 1
	}

	title := fmt.Sprintf("%s  peak: %s", label, formatRate(peak, useBits, useSI))

	// grid[row][col]: row 0 is top, row graphHeight-1 is bottom
	grid := make([][]rune, graphHeight)
	for r := range grid {
		grid[r] = make([]rune, gw)
		for c := range grid[r] {
			grid[r][c] = ' '
		}
	}

	for c, s := range samples {
		if c >= gw {
			break
		}
		frac := s / divisor
		filled := frac * float64(graphHeight) * 8
		fullRows := int(filled) / 8
		partial := int(filled) % 8

		for r := 0; r < graphHeight; r++ {
			fromBottom := graphHeight - 1 - r
			switch {
			case fromBottom < fullRows:
				grid[r][c] = '█'
			case fromBottom == fullRows && partial > 0:
				grid[r][c] = blocks[partial-1]
			}
		}
	}

	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteByte('\n')
	for r := 0; r < graphHeight; r++ {
		sb.WriteRune('│')
		sb.WriteString(string(grid[r]))
		sb.WriteByte('\n')
	}
	sb.WriteRune('└')
	sb.WriteString(strings.Repeat("─", gw))
	return sb.String()
}
