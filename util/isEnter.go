package utils

import tea "github.com/charmbracelet/bubbletea"

func IsEnter(msg tea.Msg) bool {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
		return true
	}
	return false
}
