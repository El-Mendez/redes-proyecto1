package views

import (
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
)

func HandleIncoming(client *protocol.Client) {
	// Tell the server you are ready
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

}
