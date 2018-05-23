package contact

import (
	"google.golang.org/grpc"
	"github.com/gambarini/grpcdemo/clients"
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"fmt"
)

func NewInternalContactClient() (contactClient contactpb.ContactsClient, conn *grpc.ClientConn, err error) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err = grpc.Dial(clients.InternalContactServiceName, opts...)

	if err != nil {
		return contactClient, conn, fmt.Errorf("failed to dial to Contact Service: %v", err)
	}

	contactClient = contactpb.NewContactsClient(conn)

	return contactClient, conn, nil
}
