package chat

import (
	"google.golang.org/grpc"
	"log"
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"github.com/gambarini/grpcdemo/cliutils"
	"fmt"
)

func NewInternalChatClient() (chatClient chatpb.ChatClient, conn *grpc.ClientConn, err error) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err = grpc.Dial(cliutils.InternalChatServiceName, opts...)

	if err != nil {
		return chatClient, conn, fmt.Errorf("failed to dial to Chat Service: %v", err)
	}

	chatClient = chatpb.NewChatClient(conn)

	return chatClient, conn, nil
}

func NewExternalChatClient() (chatClient chatpb.ChatClient, conn *grpc.ClientConn) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(cliutils.ExternalChatServiceDomain, opts...)

	if err != nil {
		log.Fatalf("failed to dial Chat Service: %v", err)
	}

	return chatpb.NewChatClient(conn), conn
}
