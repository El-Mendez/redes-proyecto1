package main

import (
	"fmt"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/util"
)

func main() {
	utils.InitializeLogger()
	defer utils.Logger.Sync()

	var input string
	fmt.Print("Ingresa tu JID: ")
	fmt.Scanln(&input)

	jid, _ := protocol.JIDFromString(input)

	stream, _ := protocol.MakeStream(jid.Domain)
	defer stream.Close()

	var password string
	fmt.Print("Ingresa tu contraseña: ")
	fmt.Scanln(&password)

	fmt.Println(password)

	client, _ := protocol.SignIn(&jid, stream, password)
	fmt.Println(client)
}
