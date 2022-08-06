package utils

import tea "github.com/charmbracelet/bubbletea"

func IsCtrlQ(msg tea.Msg) bool {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyCtrlQ {
		return true
	}
	return false
}
