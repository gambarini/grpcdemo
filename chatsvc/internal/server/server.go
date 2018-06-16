package server

import (

	"github.com/gambarini/grpcdemo/chatsvc/internal/repo"
	"google.golang.org/grpc"
	"github.com/gambarini/grpcdemo/chatsvc/internal/queue"
	"github.com/gambarini/grpcdemo/pb/messagepb"
	"github.com/gambarini/grpcdemo/svcutils"
	"github.com/gambarini/grpcdemo/dbutils"
	"github.com/gambarini/grpcdemo/cliutils/message"
	"github.com/gambarini/grpcdemo/pb/chatpb"
)

type ChatServer struct {
	Repository        *repo.ChatRepository
	//ContactClient     contactpb.ContactsClient
	//ContactClientConn *grpc.ClientConn
	MessageClient     messagepb.MessageClient
	MessageClientConn *grpc.ClientConn
	ChatMQ            *queue.ChatMQ
}

func (server *ChatServer) Initialize(main *svcutils.Main) error {

	db, err := dbutils.NewMongoDB(dbutils.MongoDBURL)

	if err != nil {
		return err
	}

	urls := []string{
		"amqp://rabbit:rabbit@rmq-0.rmq/vh_grpcdemo",
		"amqp://rabbit:rabbit@rmq-1.rmq/vh_grpcdemo",
	}

	chatMq, err := queue.NewChatMQ(urls)

	if err != nil {
		return err
	}

	messageClient, messageConn, err := message.NewInternalMessageClient()

	server.Repository = repo.NewChatRepository(db)
	server.MessageClient = messageClient
	server.MessageClientConn = messageConn
	server.ChatMQ = chatMq

	chatpb.RegisterChatServer(main.GRPCServer, server)

	return nil
}

func (server *ChatServer) CleanUp() {

	//server.ContactClientConn.Close()

	server.MessageClientConn.Close()

	for _, mqConn := range server.ChatMQ.MqConnections {
		mqConn.Close()
	}

	server.Repository.DB.CleanUp()
}
