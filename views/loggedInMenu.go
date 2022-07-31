package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
)

var loggedInOptions = []string{"Show all contacts", "Add a contact", "See a user details", "Send a message",
	"Send a message (group)", "Set a presence", "Send a file", "Log Out", "Delete Account"}

type LoggedInMenu struct {
	p      *tea.Program
	client *protocol.Client
}

func (m *LoggedInMenu) Start(client *protocol.Client, p *tea.Program) {
	m.client = client
	m.p = p
}

func InitialLoggedInMenu() *LoggedInMenu {
	return &LoggedInMenu{}
}

func (m *LoggedInMenu) Init() tea.Cmd {
	return nil
}

func (m *LoggedInMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *LoggedInMenu) View() string {
	return "You are logged in as " + m.client.FullJid()
}
