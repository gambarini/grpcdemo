package main

import (
	pb"github.com/gambarini/grpcdemo/pb/contact"
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 30002))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening on :30002")
	grpcServer.Serve(listener)
}

