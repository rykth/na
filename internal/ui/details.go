package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderDetails renders the extended counters section below the main stats panel.
func renderDetails(m Model, w, h int) string {
	names := m.store.Names()
	if len(names) == 0 || m.selected >= len(names) {
		return lipgloss.NewStyle().Width(w).Height(h).Render("")
	}
	iface := m.store.Get(names[m.selected])
	if iface == nil {
		return lipgloss.NewStyle().Width(w).Height(h).Render("")
	}

	sep := strings.Repeat("─", w)
	st := &iface.Stats

	content := strings.Join([]string{
		sep,
		"Extended counters:",
		statRow("Carrier", "—", formatCount(st.TxCarrier.Rate)),
		statRow("Compressed", formatCount(st.RxCompressed.Rate), formatCount(st.TxCompressed.Rate)),
		statRow("Multicast", formatCount(st.RxMulticast.Rate), "—"),
	}, "\n")

	return lipgloss.NewStyle().Width(w).Height(h).Render(content)
}
