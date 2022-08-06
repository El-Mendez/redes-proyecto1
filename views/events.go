package views

import "github.com/el-mendez/redes-proyecto1/protocol"

type LoggedInMsg struct {
	Client *protocol.Client
}

type LoggedOutMsg struct {
	Client *protocol.Client
}

type ErrorMsg struct {
	Err string
}
