package chat

import (
	"google.golang.org/grpc"
	"log"
	chatPb "github.com/gambarini/grpcdemo/pb/chat"
	"github.com/gambarini/grpcdemo/clients"
)


func NewChatClient() (chatClient chatPb.ChatClient, conn *grpc.ClientConn) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(clients.ChatDomain, opts...)

	if err != nil {
		log.Fatalf("failed to dial Chat Service: %v", err)
	}

	return chatPb.NewChatClient(conn), conn
}


