package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"sync"
)

type Device struct {
	Show   string
	Status []string
}

type GlobalState struct {
	P            *tea.Program
	Client       *protocol.Client
	FriendsMutex sync.Mutex
	//			friend		devices
	Friends map[string]map[string]*Device

	// like the alerts when a user logins
	AlertStyle lipgloss.Style

	// like private message
	TypeStyle   lipgloss.Style
	SenderStyle lipgloss.Style

	// any warning
	WarningStyle lipgloss.Style
}

var State = &GlobalState{
	AlertStyle: lipgloss.NewStyle().Faint(true),
	SenderStyle: lipgloss.NewStyle().
		Background(lipgloss.Color("86")).
		Foreground(lipgloss.Color("0")).
		Bold(true),
	TypeStyle: lipgloss.NewStyle().
		Background(lipgloss.Color("86")).
		Foreground(lipgloss.Color("9")).
		Bold(true),
	WarningStyle: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true),
}
