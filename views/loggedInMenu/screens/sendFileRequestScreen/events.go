package sendFileRequestScreen

import (
	"encoding/base64"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas/query"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"github.com/el-mendez/redes-proyecto1/views"
	"io/ioutil"
)

const CHUNK_SIZE = 320

func getData(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(content), nil
}

func sendFileRequest(username string, filename string) tea.Cmd {
	return func() tea.Msg {
		var content, id, sid string
		var seq int
		var err error

		content, err = getData(filename)
		if err != nil {
			views.State.P.Send(views.NotificationAndBack{
				Msg: views.State.AlertStyle.Render("You tried to send a file that does not exist."),
			})
		}

		var responses = make(chan *stanzas.IQ)

		// Prepare for using the channels callback
		id = stanzas.GenerateID()
		sid = stanzas.GenerateID()

		views.State.AddChannel(id, responses)

		// Send the IBB invitation
		views.State.Client.Send <- &stanzas.IQ{
			ID:   id,
			Type: "set",
			To:   username,
			Query: &query.OpenIBBQuery{
				BlockSize: 4096,
				Sid:       sid,
			},
		}

		views.State.P.Send(views.NotificationAndBack{
			Msg: views.State.AlertStyle.Render("You tried to send " + filename + " to " + username),
		})

		for seq*CHUNK_SIZE < len(content) {
			response := <-responses
			views.State.DeleteChannel(id)

			if response.Type != "result" {
				return views.Notification{
					Msg: views.State.AlertStyle.Render(filename + " could not be sent to " + username + ". :("),
				}
			}

			id = stanzas.GenerateID()
			views.State.AddChannel(id, responses)
			views.State.Client.Send <- &stanzas.IQ{
				ID:   id,
				Type: "set",
				To:   username,
				Query: &query.IBBDataQuery{
					Value:    content[seq*CHUNK_SIZE : utils.Min((seq+1)*CHUNK_SIZE, len(content))],
					Sequence: uint16(seq),
					Sid:      sid,
				},
			}
			seq++
		}

		response := <-responses
		views.State.DeleteChannel(id)
		close(responses)

		if response.Type != "result" {
			return views.Notification{
				Msg: views.State.AlertStyle.Render(filename + " could not be sent to " + username + ". :("),
			}
		}

		id = stanzas.GenerateID()
		views.State.Client.Send <- &stanzas.IQ{
			ID:    id,
			Type:  "set",
			To:    username,
			Query: &query.CloseIBBQuery{Sid: sid},
		}

		return views.NotificationAndBack{
			Msg: views.State.AlertStyle.Render("You successfully sent " + filename + " to " + username),
		}
	}
}
