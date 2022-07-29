package main

import (
	"flag"
	"fmt"
	"github.com/el-mendez/redes-proyecto1/protocol"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/util"
)

func main() {
	utils.InitializeLogger()
	defer utils.Logger.Sync()

	var account, password string

	flag.StringVar(&account, "account", "", "The JID account to log in with")
	flag.StringVar(&password, "password", "", "The matching password for the account")
	flag.Parse()

	jid, ok := protocol.JIDFromString(account)
	if !ok {
		utils.Logger.Fatal("You entered an invalid account.")
	}

	client, err := protocol.LogIn(&jid, password)
	if client != nil {
		defer client.Close()
	}
	if err != nil {
		fmt.Printf("Could not log in: %v", err)
	}

	client.Send <- &stanzas.Message{
		Type: "chat",
		To:   "mendez@alumchat.fun",
		From: client.FullJid(),
		Body: "Hola desde channels!",
	}
}
