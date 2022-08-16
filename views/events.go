package views

import "github.com/el-mendez/redes-proyecto1/protocol"

type LoggedInMsg struct {
	Client *protocol.Client
}

type LoggedOutMsg struct {
}

type ErrorMsg struct {
	Err string
}

type Notification struct {
	Msg string
}

type NotificationAndBack struct {
	Msg string
}

type FriendRequest struct {
	From string
}

type FileRequest struct {
	From string
	Sid  string
	Id   string
}
