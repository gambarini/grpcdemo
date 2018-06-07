package main

import (
	"github.com/gambarini/grpcdemo/svcutils"
	"github.com/gambarini/grpcdemo/contactsvc/internal/server"
)

func main() {

	main := svcutils.Main{
		ServerPort:     80,
		Name:           "Contact Service",
		Server: &server.ContactsServer{},
	}

	main.Run()
}


