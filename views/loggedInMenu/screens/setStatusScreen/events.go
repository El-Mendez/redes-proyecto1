package setStatusScreen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/views"
)

func setStatus(show string, message string) tea.Cmd {
	if message != "" {
		return func() tea.Msg {
			views.State.Client.Send <- &stanzas.Presence{Show: show, Status: []string{message}}
			return views.NotificationAndBack{Msg: views.State.AlertStyle.Render("You have changed your status to: " + message)}
		}
	}

	return func() tea.Msg {
		views.State.Client.Send <- &stanzas.Presence{Show: show}
		return views.NotificationAndBack{Msg: views.State.AlertStyle.Render("You have changed your status")}
	}
}
