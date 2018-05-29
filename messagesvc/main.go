package main

import (
	"github.com/gambarini/grpcdemo/svcutils"
	"github.com/gambarini/grpcdemo/dbutils"
	"github.com/gambarini/grpcdemo/messagesvc/internal/repo"
	"github.com/gambarini/grpcdemo/messagesvc/internal/server"
	"github.com/gambarini/grpcdemo/pb/messagepb"
)

func main() {

	mainServer := svcutils.MainServer{
		Initialization: initialization,
		CleanUp:        cleanUp,
		ServerPort:     30003,
		Name:           "Message Service",
	}

	mainServer.Run()
}

func initialization(mainServer *svcutils.MainServer) (err error) {

	db, err := dbutils.NewMongoDB(dbutils.MongoDBURL)

	if err != nil {
		return err
	}
	messageServer := &server.MessageServer{
		Repository:        repo.NewMessageRepository(db),
	}

	messagepb.RegisterMessageServer(mainServer.GRPCServer, messageServer)

	mainServer.Server = messageServer

	return nil
}

func cleanUp(mainServer *svcutils.MainServer) {

	chatServer := mainServer.Server.(*server.MessageServer)

	chatServer.Repository.DB.CleanUp()
}
