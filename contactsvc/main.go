package main

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"github.com/gambarini/grpcdemo/contactsvc/internal/server"
	"github.com/gambarini/grpcdemo/contactsvc/internal/db"
	"gopkg.in/mgo.v2"
	"github.com/gambarini/grpcdemo/svcutils"
	"fmt"
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

	session, err := mgo.Dial("mongodb-set-0.mongodb-service,mongodb-set-1.mongodb-service,mongodb-set-2.mongodb-service")

	if err != nil {
		return fmt.Errorf("fail to dial to mongodb cluster, %s", err)
	}

	contactsServer := &server.ContactsServer{
		DB: db.NewDB(session),
	}

	contactpb.RegisterContactsServer(mainServer.GRPCServer, contactsServer)

	mainServer.ServerObjects = contactsServer

	return nil
}

func cleanUp(mainServer *svcutils.MainServer) {

	contactsServer := mainServer.ServerObjects.(*server.ContactsServer)

	contactsServer.DB.Session.Close()
}
