package main

import (
	"github.com/gambarini/grpcdemo/svcutils"
	"github.com/gambarini/grpcdemo/contactsvc/internal/server"
)

func main() {

	main := svcutils.Main{
		ServerPort: 50051,
		Name:       "Contact Service",
		Server:     &server.ContactsServer{},
	}

	main.Run()
}
