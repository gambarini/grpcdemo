package server

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"github.com/gambarini/grpcdemo/chatsvc/internal/repo"
	"google.golang.org/grpc"
	"github.com/gambarini/grpcdemo/chatsvc/internal/queue"
	"github.com/gambarini/grpcdemo/pb/messagepb"
	"github.com/gambarini/grpcdemo/svcutils"
	"github.com/gambarini/grpcdemo/dbutils"
	"github.com/streadway/amqp"
	"fmt"
	"github.com/gambarini/grpcdemo/cliutils/contact"
	"github.com/gambarini/grpcdemo/cliutils/message"
	"github.com/gambarini/grpcdemo/pb/chatpb"
)

type ChatServer struct {
	Repository        *repo.ChatRepository
	ContactClient     contactpb.ContactsClient
	ContactClientConn *grpc.ClientConn
	MessageClient     messagepb.MessageClient
	MessageClientConn *grpc.ClientConn
	ChatMQ            *queue.ChatMQ
}

func (server *ChatServer) Initialize(main *svcutils.Main) error {

	db, err := dbutils.NewMongoDB(dbutils.MongoDBURL)

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

	messageClient, messageConn, err := message.NewInternalMessageClient()

	server.Repository = repo.NewChatRepository(db)
	server.ContactClient = contactClient
	server.ContactClientConn = conn
	server.MessageClient = messageClient
	server.MessageClientConn = messageConn
	server.ChatMQ = queue.NewChatMQ(mqConnection)

	chatpb.RegisterChatServer(main.GRPCServer, server)

	return nil
}

func (server *ChatServer) CleanUp() {

	server.ContactClientConn.Close()

	server.MessageClientConn.Close()

	server.ChatMQ.MqConnection.Close()

	server.Repository.DB.CleanUp()
}
