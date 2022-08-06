package LoginScreen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/views"
)

type LoggedInMsg struct {
	Err error
}

func logIn(jid protocol.JID, password string) tea.Cmd {
	return func() tea.Msg {
		c, err := protocol.LogIn(&jid, password)
		if err != nil {
			c.Close()
		} else {
			views.State.Client = c
		}

		return LoggedInMsg{err}
	}
}
