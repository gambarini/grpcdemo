package message

import (
	"google.golang.org/grpc"
	"github.com/gambarini/grpcdemo/cliutils"
	"fmt"

	"github.com/gambarini/grpcdemo/pb/messagepb"
	"log"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
)

func NewInternalMessageClient() (contactClient messagepb.MessageClient, conn *grpc.ClientConn, err error) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err = grpc.Dial(cliutils.InternalMessageServiceName, opts...)

	if err != nil {
		return contactClient, conn, fmt.Errorf("failed to dial to Message Service: %v", err)
	}

	contactClient = messagepb.NewMessageClient(conn)

	return contactClient, conn, nil
}

func NewExternalMessageClient() (messageClient messagepb.MessageClient, conn *grpc.ClientConn) {

	creds := credentials.NewTLS(&tls.Config{ InsecureSkipVerify: true})

	var opts []grpc.DialOption

	//opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTransportCredentials(creds))
	//opts = append(opts, grpc.WithStreamInterceptor(cliutils.StreamClientInterceptor))
	//opts = append(opts, grpc.WithUnaryInterceptor(cliutils.UnaryClientInterceptor))

	conn, err := grpc.Dial(cliutils.ExternalDomainMessage, opts...)

	if err != nil {
		log.Fatalf("failed to dial Message Service: %v", err)
	}

	return messagepb.NewMessageClient(conn), conn
}
