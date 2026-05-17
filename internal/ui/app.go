package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rykth/na/internal/collector"
	"github.com/rykth/na/internal/config"
	"github.com/rykth/na/internal/model"
)

// collectorMsg carries a snapshot from one collector.Read() call.
type collectorMsg struct {
	snapshots []collector.RawIfStats
	at        time.Time
	err       error
}

// Model is the Bubbletea root model.
type Model struct {
	width, height int
	store         *model.Store
	collector     collector.Collector
	selected      int // index into store.Names()
	interval      time.Duration
	cfg           *config.Config
	showHelp      bool
	useBits       bool
	useSI         bool
	err           error // last collector error, shown in the status bar
}

// New constructs the root UI model.
func New(col collector.Collector, store *model.Store, cfg *config.Config) Model {
	return Model{
		store:     store,
		collector: col,
		interval:  time.Duration(cfg.Interval * float64(time.Second)),
		cfg:       cfg,
	}
}

// Init fires the first collector read immediately.
func (m Model) Init() tea.Cmd {
	return m.doCollect()
}

// doCollect returns a Cmd that reads one snapshot from the collector.
func (m Model) doCollect() tea.Cmd {
	return func() tea.Msg {
		snaps, err := m.collector.Read()
		return collectorMsg{snapshots: snaps, at: time.Now(), err: err}
	}
}

// scheduleNext returns a Cmd that waits one interval then runs another collect.
func (m Model) scheduleNext() tea.Cmd {
	return tea.Tick(m.interval, func(_ time.Time) tea.Msg {
		return m.doCollect()()
	})
}

// Update is the Bubbletea event handler.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case collectorMsg:
		m.err = msg.err
		if msg.err == nil && len(msg.snapshots) > 0 {
			m.store.Update(msg.snapshots, msg.at)
			// auto-select first interface on startup
			if m.selected == 0 && len(m.store.Names()) > 0 {
				m.selected = 0
			}
		}
		return m, m.scheduleNext()

	case tea.KeyMsg:
		// dismiss help on any key
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?", "h":
			m.showHelp = true
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			names := m.store.Names()
			if m.selected < len(names)-1 {
				m.selected++
			}
		}
		return m, nil
	}
	return m, nil
}

// View assembles the full-screen layout.
func (m Model) View() string {
	if m.width == 0 {
		return "loading…"
	}
	if m.showHelp {
		return renderHelp(m)
	}

	bodyH := m.height - 2 // subtract header and status bar rows
	rightW := m.width - listWidth - 1

	left := renderIfList(m, bodyH)
	right := renderStats(m, rightW, bodyH)
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	return lipgloss.JoinVertical(lipgloss.Left,
		renderHeader(m),
		body,
		renderStatusBar(m),
	)
}

func renderHeader(m Model) string {
	title := "na — network analyzer"
	ts := time.Now().Format("2006-01-02 15:04:05")
	pad := max(0, m.width-len(title)-len(ts)-2) // -2 for padding(0,1)
	content := title + fmt.Sprintf("%*s%s", pad, "", ts)
	return styleHeader.Width(m.width).Render(content)
}

func renderStatusBar(m Model) string {
	if m.err != nil {
		return styleStatus.Width(m.width).Render(fmt.Sprintf("error: %v", m.err))
	}
	interval := fmt.Sprintf("interval: %.1fs", m.interval.Seconds())
	hint := "q quit  ? help  ↑↓/jk select"
	pad := max(1, m.width-len(interval)-len(hint))
	content := interval + fmt.Sprintf("%*s%s", pad, "", hint)
	return styleStatus.Width(m.width).Render(content)
}
