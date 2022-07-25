package main

import (
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/util"
)

func main() {
	utils.InitializeLogger()
	defer utils.Logger.Sync()

	stream, _ := protocol.MakeStream("alumchat.fun")
	defer stream.Close()
}
