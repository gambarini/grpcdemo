package main

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"github.com/gambarini/grpcdemo/contactsvc/internal/server"
	"github.com/gambarini/grpcdemo/svcutils"
	"github.com/gambarini/grpcdemo/dbutils"
	"github.com/gambarini/grpcdemo/contactsvc/internal/repo"
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

	db, err := dbutils.NewMongoDB(dbutils.MongoDBURL)

	if err != nil {
		return err
	}

	contactsServer := &server.ContactsServer{
		ContactRepository: repo.NewContactRepository(db),
	}

	contactpb.RegisterContactsServer(mainServer.GRPCServer, contactsServer)

	mainServer.Server = contactsServer

	return nil
}

func cleanUp(mainServer *svcutils.MainServer) {

	contactsServer := mainServer.Server.(*server.ContactsServer)

	contactsServer.ContactRepository.DB.CleanUp()
}
