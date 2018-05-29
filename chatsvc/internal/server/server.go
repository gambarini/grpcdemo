package server

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"github.com/gambarini/grpcdemo/chatsvc/internal/repo"
	"google.golang.org/grpc"
	"github.com/gambarini/grpcdemo/chatsvc/internal/queue"
	"github.com/gambarini/grpcdemo/pb/messagepb"
)

type ChatServer struct {
	Repository        *repo.ChatRepository
	ContactClient     contactpb.ContactsClient
	ContactClientConn *grpc.ClientConn
	MessageClient     messagepb.MessageClient
	MessageClientConn *grpc.ClientConn
	ChatMQ            *queue.ChatMQ
}


