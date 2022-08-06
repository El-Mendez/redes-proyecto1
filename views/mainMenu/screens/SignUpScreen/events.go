package SignUpScreen

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/views"
)

func signUp(jid protocol.JID, password string) tea.Cmd {
	return func() tea.Msg {
		c, err := protocol.SignUp(&jid, password)
		if err != nil {
			if c != nil {
				c.Close()
			}
			return views.ErrorMsg{Err: views.State.WarningStyle.Render(fmt.Sprintf("Could not sign up: %v", err))}
		}

		return views.LoggedInMsg{Client: c}
	}
}
