package server

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"github.com/gambarini/grpcdemo/chatsvc/internal/db"
	"google.golang.org/grpc"
	"github.com/gambarini/grpcdemo/chatsvc/internal/queue"
)

type ChatServer struct {
	DB                *db.DB
	ContactClient     contactpb.ContactsClient
	ContactClientConn *grpc.ClientConn
	ChatMQ            *queue.ChatMQ
}


