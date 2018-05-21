package contact

import (
	"google.golang.org/grpc"
	"log"
	"github.com/gambarini/grpcdemo/clients"
	contactPb "github.com/gambarini/grpcdemo/pb/contact"
)

func NewInternalContactClient() (contactClient contactPb.ContactsClient, conn *grpc.ClientConn) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(clients.InternalContactServiceName, opts...)

	if err != nil {
		log.Fatalf("failed to dial Contact Service: %v", err)
	}

	return contactPb.NewContactsClient(conn), conn
}
