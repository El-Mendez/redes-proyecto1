package loggedInMenu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/views"
)

func logout() tea.Msg {
	views.State.Client.Send <- &stanzas.Presence{Type: "unavailable"}
	views.State.Client.Close()
	return views.LoggedOutMsg{}
}

func deleteAccount() tea.Msg {
	views.State.Client.Send <- &stanzas.Presence{Type: "unavailable"}
	views.State.Client.DeleteAccount()
	views.State.Client.Close()
	return views.LoggedOutMsg{}
}
