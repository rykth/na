package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rykth/na/internal/collector"
	_ "github.com/rykth/na/internal/collector/linux"
	"github.com/rykth/na/internal/config"
	"github.com/rykth/na/internal/model"
	"github.com/rykth/na/internal/ui"
)

func main() {
	cfg := config.Parse()

	col, err := collector.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "na: %v\n", err)
		os.Exit(1)
	}
	defer col.Close()

	store := model.NewStore(cfg.HistorySize)
	m := ui.New(col, store, cfg)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "na: %v\n", err)
		os.Exit(1)
	}
}
