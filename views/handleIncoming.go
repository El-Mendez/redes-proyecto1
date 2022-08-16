package views

import (
	"encoding/base64"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas/query"
	"os"
)

func HandleIncoming(client *protocol.Client) {
	// Tell the server you are ready
	client.Send <- &stanzas.IQ{
		ID:    stanzas.GenerateID(),
		Type:  "get",
		To:    client.BaseJid(),
		From:  client.FullJid(),
		Query: &query.RosterQuery{},
	}
	client.Send <- &stanzas.Presence{}

	for stanza := range client.Receive {
		switch s := stanza.(type) {
		case *stanzas.Message:
			handleIncomingMessages(s)
		case *stanzas.Presence:
			handleIncomingPresences(s)
		case *stanzas.IQ:
			handleIncomingIQ(s)
		}
	}
}

func handleIncomingMessages(msg *stanzas.Message) {
	if msg.From == "" || msg.Body == "" {
		return
	}

	switch msg.Type {
	case "groupchat":
		from, _ := protocol.JIDFromString(msg.From)
		State.P.Send(Notification{
			Msg: State.SenderStyle.Render(from.DeviceName) +
				State.TypeStyle.Render(" through "+from.BaseJid()) +
				" " + msg.Body,
		})
	default:
		State.P.Send(Notification{
			Msg: State.SenderStyle.Render(msg.From) +
				State.TypeStyle.Render(" to you") +
				" " + msg.Body,
		})
	}

}

func handleIncomingPresences(presence *stanzas.Presence) {
	// Handle subscription
	switch presence.Type {
	case "subscribed":
		State.P.Send(Notification{State.AlertStyle.Render(presence.From + " accepted your friend request!")})
		return
	case "unsubscribed":
		State.P.Send(Notification{State.AlertStyle.Render(presence.From + " has stopped being your friend.")})
		return
	case "subscribe":
		State.P.Send(FriendRequest{From: presence.From})
		return
	case "unavailable":
		jid, ok := protocol.JIDFromString(presence.From)
		if !ok || jid.DeviceName == "" {
			return
		}

		State.FriendsMutex.Lock()
		defer State.FriendsMutex.Unlock()

		if friend, ok := State.Friends[jid.BaseJid()]; ok {
			if _, ok := friend[jid.DeviceName]; ok {
				delete(friend, jid.DeviceName)
			}
		}

		State.P.Send(Notification{State.AlertStyle.Render(presence.From + " has disconnected.")})
		return
	}

	// Regular status change
	switch presence.Show {
	case "chat":
		State.P.Send(Notification{State.AlertStyle.Render(presence.From + " is Available")})
	case "away":
		State.P.Send(Notification{State.AlertStyle.Render(presence.From + " is Away")})
	case "xa":
		State.P.Send(Notification{State.AlertStyle.Render(presence.From + " is Extended Away")})
	case "dnd":
		State.P.Send(Notification{State.AlertStyle.Render(presence.From + " is Busy")})
	}

	if presence.Status != nil && len(presence.Status) > 0 {
		data := ""
		for _, state := range presence.Status {
			data += "\n\t" + state
		}
		State.P.Send(Notification{State.AlertStyle.Render(presence.From + " changed his state to: " + data)})
	}

	jid, ok := protocol.JIDFromString(presence.From)
	if !ok || jid.DeviceName == "" {
		return
	}

	State.FriendsMutex.Lock()
	defer State.FriendsMutex.Unlock()

	if friend, ok := State.Friends[jid.BaseJid()]; ok {
		if _, ok := friend[jid.DeviceName]; !ok {
			State.P.Send(Notification{State.AlertStyle.Render(presence.From) + " has connected. "})
		}
		if presence.Show == "" {
			presence.Show = "chat"
		}
		friend[jid.DeviceName] = &Device{Show: presence.Show, Status: presence.Status}
	}
}

func handleIncomingIQ(iq *stanzas.IQ) {
	State.ChannelsMutex.Lock()
	defer State.ChannelsMutex.Unlock()

	if channel, ok := State.Channels[iq.ID]; ok {
		channel <- iq
	}

	switch q := iq.Query.(type) {
	case *query.RosterQuery:
		if iq.Type == "result" || iq.Type == "set" {
			if q.RosterItems == nil {
				return
			}
			State.FriendsMutex.Lock()
			defer State.FriendsMutex.Unlock()
			for _, friend := range q.RosterItems {
				if _, ok := State.Friends[friend.Jid]; !ok {
					State.Friends[friend.Jid] = make(map[string]*Device)
				}
			}
		}
	case *query.OpenIBBQuery:
		State.P.Send(FileRequest{From: iq.From, Sid: q.Sid, Id: iq.ID})
	case *query.IBBDataQuery:
		State.FileMutex.Lock()
		if transaction, ok := State.FileTransactions[q.Sid]; ok {
			transaction.Content.WriteString(q.Value)
		}
		State.FileMutex.Unlock()

		State.Client.Send <- &stanzas.IQ{ID: iq.ID, Type: "result", To: iq.From}
	case *query.CloseIBBQuery:
		State.FileMutex.Lock()
		data, ok := State.FileTransactions[q.Sid]
		if ok {
			delete(State.FileTransactions, q.Sid)
		}
		State.FileMutex.Unlock()

		if !ok {
			return
		}

		if contents, err := base64.StdEncoding.DecodeString(data.Content.String()); err == nil {
			if err := os.WriteFile(data.Filename, contents, 0777); err == nil {
				State.P.Send(Notification{State.AlertStyle.Render("You successfully got " + data.Filename)})
			} else {
				State.P.Send(Notification{State.AlertStyle.Render(data.Filename + " could not be saved.")})
			}
		}
	}
}
