package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"strings"
	"sync"
)

type Device struct {
	Show   string
	Status []string
}

type FileStatus struct {
	Filename   string
	Content    strings.Builder
	CurrentSeq uint16
}

type GlobalState struct {
	P            *tea.Program
	Client       *protocol.Client
	FriendsMutex sync.Mutex
	//			friend		devices
	Friends map[string]map[string]*Device

	// Handlers through channels for simplification
	ChannelsMutex sync.Mutex
	Channels      map[string]chan<- *stanzas.IQ

	// Handlers through channels for simplification
	FileMutex        sync.Mutex
	FileTransactions map[string]*FileStatus

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

func (s *GlobalState) AddChannel(id string, channel chan<- *stanzas.IQ) {
	s.ChannelsMutex.Lock()
	defer s.ChannelsMutex.Unlock()

	s.Channels[id] = channel
}

func (s *GlobalState) DeleteChannel(id string) {
	s.ChannelsMutex.Lock()
	defer s.ChannelsMutex.Unlock()

	delete(s.Channels, id)
}
