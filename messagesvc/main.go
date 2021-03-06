package main

import (
	"github.com/gambarini/grpcdemo/svcutils"
	"github.com/gambarini/grpcdemo/messagesvc/internal/server"
)

func main() {

	main := svcutils.Main{
		ServerPort: 50051,
		Name:       "Message Service",
		Server:     &server.MessageServer{},
	}

	main.Run()
}
