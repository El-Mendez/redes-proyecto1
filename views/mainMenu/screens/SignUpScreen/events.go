package SignUpScreen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/views"
)

type SignedUpMsg struct {
	Err error
}

func signUp(jid protocol.JID, password string) tea.Cmd {
	return func() tea.Msg {
		c, err := protocol.SignUp(&jid, password)
		if err != nil {
			c.Close()
		} else {
			views.State.Client = c
		}

		return SignedUpMsg{err}
	}
}
