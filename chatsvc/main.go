package main

import (
	"github.com/gambarini/grpcdemo/chatsvc/internal/server"
	"github.com/gambarini/grpcdemo/svcutils"
)

func main() {

	main := svcutils.Main{
		ServerPort: 30001,
		Name:       "Chat Service",
		Server:     &server.ChatServer{},
	}

	main.Run()
}
