package main

import (
	pb "github.com/gambarini/grpcdemo/pb/contact"
	"google.golang.org/grpc"
	"net"
	"fmt"
	"log"
	"github.com/gambarini/grpcdemo/servers/contact/internal/server"
	"github.com/gambarini/grpcdemo/servers/contact/internal/db"
)

func main() {

	//var opts []grpc.ServerOption

	grpcServer := grpc.NewServer()

	pb.RegisterContactsServer(grpcServer, &server.ContactsServer{
		DB: db.NewDB(),
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 30002))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening on port 30002")
	grpcServer.Serve(listener)
}
