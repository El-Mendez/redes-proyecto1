package loggedInMenu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
)

func (m *LoggedInMenu) sendMessage(client *protocol.Client, to string, content string) tea.Cmd {
	return func() tea.Msg {
		client.Send <- &stanzas.Message{
			Type: "chat",
			To:   to,
			From: client.FullJid(),
			Body: content,
		}

		return notification{
			text: m.senderStyle.Render("You to ") +
				m.typeStyle.Render(to) +
				": " +
				content,
		}
	}
}

func (m *LoggedInMenu) logOut(client *protocol.Client) tea.Cmd {
	return func() tea.Msg {
		client.Send <- &stanzas.Presence{Type: "unavailable"}
		client.Close()
		return LogOutResult{}
	}
}

func (m *LoggedInMenu) deleteAccount(client *protocol.Client) tea.Cmd {
	return func() tea.Msg {
		client.Send <- &stanzas.Presence{Type: "unavailable"}
		client.DeleteAccount()
		client.Close()
		return LogOutResult{}
	}
}
