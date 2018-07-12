package contact

import (
	"google.golang.org/grpc"
	"github.com/gambarini/grpcdemo/cliutils"
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"fmt"
	"google.golang.org/grpc/keepalive"
	"time"
	"log"
)

func NewInternalContactClient() (contactClient contactpb.ContactsClient, conn *grpc.ClientConn, err error) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err = grpc.Dial(cliutils.InternalContactServiceName, opts...)

	if err != nil {
		return contactClient, conn, fmt.Errorf("failed to dial to Contact Service: %v", err)
	}

	contactClient = contactpb.NewContactsClient(conn)

	return contactClient, conn, nil
}

func NewExternalContactClient() (contactClient contactpb.ContactsClient, conn *grpc.ClientConn) {

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

	conn, err := grpc.Dial(cliutils.ExternalDomainContact, opts...)

	if err != nil {
		log.Fatalf("Failed to dial Contact Service: %v", err)
	}

	return contactpb.NewContactsClient(conn), conn
}