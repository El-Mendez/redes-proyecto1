package utils

import "github.com/charmbracelet/lipgloss"

func MenuOption(value string, selected bool, style lipgloss.Style) string {
	if selected {
		return style.Render("[✓] " + value)
	}
	return "[ ] " + value
}
