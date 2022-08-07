package views

import (
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas/query"
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

	switch presence.Type {
	case "subscribed":
		State.P.Send(Notification{State.AlertStyle.Render(presence.From + " accepted your friend request!")})
	case "unsubscribed":
		State.P.Send(Notification{State.AlertStyle.Render(presence.From + " has stopped being your friend.")})
	case "subscribe":
		State.P.Send(FriendRequest{From: presence.From})
	}
}

func handleIncomingIQ(iq *stanzas.IQ) {
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
	}
}
