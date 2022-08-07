package sendFriendRequestScreen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/views"
)

func sendFriendRequest(username string) tea.Cmd {
	return func() tea.Msg {
		views.State.Client.Send <- &stanzas.Presence{
			To:   username,
			Type: "subscribe",
		}
		return views.NotificationAndBack{
			Msg: views.State.AlertStyle.Render("You sent a friend request to " + username),
		}
	}
}
