package gotFriendRequest

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/views"
)

func acceptFriendRequest(username string) tea.Cmd {
	return func() tea.Msg {
		views.State.Client.Send <- &stanzas.Presence{To: username, Type: "subscribed"}
		return views.NotificationAndBack{
			Msg: views.State.AlertStyle.Render("You are now friends with " + username),
		}
	}
}

func rejectFriendshipRequest(username string) tea.Cmd {
	return func() tea.Msg {
		views.State.Client.Send <- &stanzas.Presence{To: username, Type: "unsubscribed"}
		return views.NotificationAndBack{
			Msg: views.State.AlertStyle.Render("You are not friends with " + username),
		}
	}
}
