package joinGroupScreen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/views"
)

func joinGroup(group string, alias string) tea.Cmd {
	return func() tea.Msg {
		views.State.Client.Send <- &stanzas.Presence{
			To: group + "/" + alias,
		}
		return views.NotificationAndBack{
			Msg: views.State.AlertStyle.Render("You tried to join the group " + group + " as " + alias),
		}
	}
}
