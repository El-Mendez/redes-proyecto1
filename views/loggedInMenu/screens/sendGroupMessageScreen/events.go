package sendGroupMessageScreen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/views"
)

func sendGroupMessage(to string, content string) tea.Cmd {
	return func() tea.Msg {
		views.State.Client.Send <- &stanzas.Message{
			Type: "groupchat",
			To:   to,
			Body: content,
		}
		return nil
		//return views.Notification{
		//	Msg: views.State.SenderStyle.Render("You to group") + views.State.TypeStyle.Render(to) + " " + content,
		//}
	}
}
