package sendFileRequestScreen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas/query"
	"github.com/el-mendez/redes-proyecto1/views"
)

func sendFileRequest(username string, filename string) tea.Cmd {
	return func() tea.Msg {
		views.State.Client.Send <- &stanzas.IQ{
			ID:   stanzas.GenerateID(),
			Type: "set",
			To:   username,
			Query: &query.OpenIBBQuery{
				BlockSize: 4096,
				Sid:       stanzas.GenerateID(),
			},
		}
		return views.NotificationAndBack{
			Msg: views.State.AlertStyle.Render("You sent " + filename + " to " + username),
		}
	}
}
