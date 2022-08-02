package mainMenu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"time"
)

func (m *MainMenu) signup(jid protocol.JID, password string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		c, err := protocol.SignUp(&jid, password)
		return LoginResult{c, err}
	}
}

func (m *MainMenu) login(jid protocol.JID, password string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		c, err := protocol.LogIn(&jid, password)
		return LoginResult{c, err}
	}
}
