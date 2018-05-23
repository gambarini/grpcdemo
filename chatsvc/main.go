package main

import (
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"github.com/gambarini/grpcdemo/chatsvc/internal/server"
	"github.com/gambarini/grpcdemo/chatsvc/internal/db"
	"github.com/gambarini/grpcdemo/clients/contact"
	"github.com/gambarini/grpcdemo/svcutils"
)

func main() {

	mainServer := svcutils.MainServer{
		Initialization: initialization,
		CleanUp:        cleanUp,
		ServerPort:     30001,
		Name:           "Chat Service",
	}

	mainServer.Run()
}

func initialization(mainServer *svcutils.MainServer) (err error) {

	contactClient, conn, err := contact.NewInternalContactClient()

	if err != nil {
		return err
	}

	chatServer := &server.ChatServer{
		DB:                db.NewDB(),
		ContactClient:     contactClient,
		ContactClientConn: conn,
	}

	chatpb.RegisterChatServer(mainServer.GRPCServer, chatServer)

	mainServer.ServerObjects = chatServer

	return nil
}

func cleanUp(mainServer *svcutils.MainServer) {

	chatServer := mainServer.ServerObjects.(*server.ChatServer)

	chatServer.ContactClientConn.Close()
}
