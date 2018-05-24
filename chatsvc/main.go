package main

import (
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"github.com/gambarini/grpcdemo/chatsvc/internal/server"
	"github.com/gambarini/grpcdemo/chatsvc/internal/db"
	"github.com/gambarini/grpcdemo/clients/contact"
	"github.com/gambarini/grpcdemo/svcutils"
	"github.com/gambarini/grpcdemo/dbutils"
	"github.com/streadway/amqp"
	"fmt"
	"github.com/gambarini/grpcdemo/chatsvc/internal/queue"
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

	session, err := dbutils.DialMongoDB()

	if err != nil {
		return err
	}

	mqConnection, err := amqp.Dial("amqp://rabbit:rabbit@172.17.0.7")

	if err != nil {
		return fmt.Errorf("failed to dial to rabbitmq cluster, %s", err)
	}

	contactClient, conn, err := contact.NewInternalContactClient()

	if err != nil {
		return err
	}

	chatServer := &server.ChatServer{
		DB:                db.NewDB(session),
		ContactClient:     contactClient,
		ContactClientConn: conn,
		ChatMQ:            queue.NewChatMQ(mqConnection),
	}

	chatpb.RegisterChatServer(mainServer.GRPCServer, chatServer)

	mainServer.Server = chatServer

	return nil
}

func cleanUp(mainServer *svcutils.MainServer) {

	chatServer := mainServer.Server.(*server.ChatServer)

	chatServer.ContactClientConn.Close()

	chatServer.ChatMQ.MqConnection.Close()
}
