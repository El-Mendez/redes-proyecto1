package main

import (
	"flag"
	"github.com/el-mendez/redes-proyecto1/protocol"
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

	stream, _ := protocol.MakeStream(jid.Domain)
	defer stream.Close()

	_, _ = protocol.SignIn(&jid, stream, password)
}
