package main

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"github.com/gambarini/grpcdemo/contactsvc/internal/server"
	"github.com/gambarini/grpcdemo/contactsvc/internal/db"
	"github.com/gambarini/grpcdemo/svcutils"
	"github.com/gambarini/grpcdemo/dbutils"
)

func main() {

	mainServer := svcutils.MainServer{
		Initialization: initialization,
		CleanUp:        cleanUp,
		ServerPort:     30002,
		Name:           "Contact Service",
	}

	mainServer.Run()
}

func initialization(mainServer *svcutils.MainServer) (err error) {

	session, err := dbutils.DialMongoDB()

	if err != nil {
		return err
	}

	contactsServer := &server.ContactsServer{
		DB: db.NewDB(session),
	}

	contactpb.RegisterContactsServer(mainServer.GRPCServer, contactsServer)

	mainServer.Server = contactsServer

	return nil
}

func cleanUp(mainServer *svcutils.MainServer) {

	contactsServer := mainServer.Server.(*server.ContactsServer)

	contactsServer.DB.Session.Close()
}
