package views

import tea "github.com/charmbracelet/bubbletea"

type Screen interface {
	tea.Model

	Focus()
	Blur()
}
