package sendMessageScreen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/views"
)

func sendMessage(to string, content string) tea.Cmd {
	return func() tea.Msg {
		views.State.Client.Send <- &stanzas.Message{
			Type: "chat",
			To:   to,
			Body: content,
		}
		return views.Notification{
			Msg: views.State.SenderStyle.Render("You to ") + views.State.TypeStyle.Render(to) + " " + content,
		}
	}
}
