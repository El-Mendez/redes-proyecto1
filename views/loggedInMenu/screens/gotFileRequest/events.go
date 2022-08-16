package gotFileRequest

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/views"
)

func acceptFileRequest(from string, id string, sid string, filename string) tea.Cmd {
	return func() tea.Msg {
		views.State.FileMutex.Lock()
		views.State.FileTransactions[sid] = &views.FileStatus{Filename: filename}
		views.State.FileMutex.Unlock()

		views.State.Client.Send <- &stanzas.IQ{ID: id, Type: "result", To: from}
		return views.NotificationAndBack{
			Msg: views.State.AlertStyle.Render("You accepted a file from " + from + " with the name " + filename),
		}
	}
}

func rejectFileRequest(from string, id string) tea.Cmd {
	return func() tea.Msg {
		views.State.Client.Send <- &stanzas.IQ{ID: id, Type: "error", To: from}
		return views.NotificationAndBack{
			Msg: views.State.AlertStyle.Render("You did not accept a file from " + from),
		}
	}
}
