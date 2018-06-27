package chat

import (
	"google.golang.org/grpc"
	"log"
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"github.com/gambarini/grpcdemo/cliutils"
	"fmt"
	"google.golang.org/grpc/keepalive"
	"time"
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

	//creds := credentials.NewTLS(&tls.Config{ InsecureSkipVerify: true})

	kap := keepalive.ClientParameters{
		PermitWithoutStream: true,
		Time: time.Minute * 5,
		Timeout: time.Minute,
	}

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	//opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithKeepaliveParams(kap))
	//opts = append(opts, grpc.WithStreamInterceptor(cliutils.StreamClientInterceptor))
	//opts = append(opts, grpc.WithUnaryInterceptor(cliutils.UnaryClientInterceptor))

	conn, err := grpc.Dial(cliutils.ExternalDomainChat, opts...)

	if err != nil {
		log.Fatalf("Failed to dial Chat Service: %v", err)
	}

	return chatpb.NewChatClient(conn), conn
}


