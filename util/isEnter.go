package utils

import tea "github.com/charmbracelet/bubbletea"

// IsEnter detects a tea event and returns true if the event is an Enter key press
func IsEnter(msg tea.Msg) bool {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
		return true
	}
	return false
}
