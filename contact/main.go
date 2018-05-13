package main

import (
	"github.com/gambarini/grpcdemo/contact/pb"
	"google.golang.org/grpc"
	"net"
	"fmt"
	"log"
)

type contactsServer struct {}

func (server *contactsServer) StoreContacts(stream pb.Contacts_StoreContactsServer) error {

	return nil
}

func (server *contactsServer) ListContacts(contact *pb.Contact, stream pb.Contacts_ListContactsServer) error {

	return nil
}

func main() {

	//var opts []grpc.ServerOption

	grpcServer := grpc.NewServer()

	pb.RegisterContactsServer(grpcServer, &contactsServer{})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening on :9000")
	grpcServer.Serve(listener)
}