package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const colW = 16 // width of each RX/TX value column

// renderStats renders the right-side stats panel for the selected interface.
func renderStats(m Model, w, h int) string {
	names := m.store.Names()
	if len(names) == 0 || m.selected >= len(names) {
		return lipgloss.NewStyle().Width(w).Height(h).Render(
			lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, "No interface selected"),
		)
	}

	name := names[m.selected]
	iface := m.store.Get(name)
	if iface == nil {
		return lipgloss.NewStyle().Width(w).Height(h).Render(
			lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, "No interface selected"),
		)
	}

	sep := strings.Repeat("─", w)

	// header line
	speed := "unknown"
	if iface.Speed > 0 {
		speed = fmt.Sprintf("%d Mbps", iface.Speed)
	}
	header := fmt.Sprintf("Interface: %s   MAC: %s   MTU: %d   Speed: %s   [%s]",
		iface.Name, iface.MAC, iface.MTU, speed, strings.ToUpper(iface.OperState))

	// column header
	colHdr := fmt.Sprintf("  %-12s %*s  %*s", "", colW, "RX", colW, "TX")

	// rate rows
	rx := func(r float64) string { return formatRate(r, m.useBits, m.useSI) }
	st := &iface.Stats

	tableRows := strings.Join([]string{
		statRow("Bytes", rx(st.RxBytes.Rate), rx(st.TxBytes.Rate)),
		statRow("Packets", formatPPS(st.RxPackets.Rate), formatPPS(st.TxPackets.Rate)),
		statRow("Errors", formatCount(st.RxErrors.Rate), formatCount(st.TxErrors.Rate)),
		statRow("Dropped", formatCount(st.RxDropped.Rate), formatCount(st.TxDropped.Rate)),
		statRow("FIFO", formatCount(st.RxFifo.Rate), formatCount(st.TxFifo.Rate)),
		fmt.Sprintf("  %-12s %*s        Collisions: %s",
			"Frame", colW, formatCount(st.RxFrame.Rate), formatCount(st.TxColls.Rate)),
	}, "\n")

	// totals
	totals := fmt.Sprintf("Totals since start:\n  RX: %-12s  (%s packets)\n  TX: %-12s  (%s packets)",
		formatBytes(st.RxBytes.Total, m.useSI), formatWithCommas(st.RxPackets.Total),
		formatBytes(st.TxBytes.Total, m.useSI), formatWithCommas(st.TxPackets.Total),
	)

	content := strings.Join([]string{header, sep, colHdr, tableRows, sep, totals}, "\n")
	return lipgloss.NewStyle().Width(w).Height(h).Render(content)
}

func statRow(label, rx, tx string) string {
	return fmt.Sprintf("  %-12s %*s  %*s", label, colW, rx, colW, tx)
}

func formatPPS(rate float64) string {
	return formatWithCommas(uint64(rate)) + " pps"
}

func formatCount(rate float64) string {
	if rate < 1 {
		return "0"
	}
	return fmt.Sprintf("%d", uint64(rate))
}

func formatWithCommas(n uint64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return b.String()
}
