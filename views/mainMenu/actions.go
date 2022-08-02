package mainMenu

import "github.com/el-mendez/redes-proyecto1/protocol"

type LoginResult struct {
	Client *protocol.Client
	Err    error
}
