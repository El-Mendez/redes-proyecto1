package utils

import tea "github.com/charmbracelet/bubbletea"

// IsCtrlQ detects a tea event and returns true if the event is a Ctrl+Q press
func IsCtrlQ(msg tea.Msg) bool {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyCtrlQ {
		return true
	}
	return false
}
