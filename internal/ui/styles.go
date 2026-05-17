package ui

import "github.com/charmbracelet/lipgloss"

const listWidth = 30 // columns reserved for the left interface-list panel

var (
	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)

	styleStatus = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	styleSelected = lipgloss.NewStyle().Reverse(true)

	styleUp   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	styleDown = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red
	styleDim  = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // dim grey
)
