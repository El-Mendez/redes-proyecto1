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

	State.P.Send(Notification{
		Msg: State.SenderStyle.Render(msg.From) +
			State.TypeStyle.Render(" to you") +
			" " + msg.Body,
	})
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

//func handleIncoming(client *protocol.Client, p *tea.Program, m *LoggedInMenu) {
//	client.Send <- &stanzas.Presence{}
//	for stanza := range client.Receive {
//		switch stanza := stanza.(type) {
//		case *stanzas.Message:
//			if stanza.From != "" && stanza.Body != "" {
//				p.Send(notification{
//					text: m.senderStyle.Render(stanza.From) +
//						m.typeStyle.Render(" to you") +
//						": " + stanza.Body,
//				})
//			}
//		case *stanzas.Presence:
//			switch stanza.Type {
//			case "subscribed":
//				p.Send(notification{m.alertStyle.Render(stanza.From + " accepted your friend request!")})
//			case "unsubscribed":
//				p.Send(notification{m.alertStyle.Render(stanza.From + " has stopped being friend.")})
//			case "subscribe":
//				p.Send(friendRequest{stanza.From})
//			}
//		}
//	}
//}
