package server

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"github.com/gambarini/grpcdemo/chatsvc/internal/repo"
	"google.golang.org/grpc"
	"github.com/gambarini/grpcdemo/chatsvc/internal/queue"
)

type ChatServer struct {
	Repository        *repo.ChatRepository
	ContactClient     contactpb.ContactsClient
	ContactClientConn *grpc.ClientConn
	ChatMQ            *queue.ChatMQ
}


