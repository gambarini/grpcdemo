package main

import (
	pb "github.com/gambarini/grpcdemo/pb/chat"
	"google.golang.org/grpc"
	"fmt"
	"log"
	"net"
	"github.com/gambarini/grpcdemo/servers/chat/internal/server"
	"github.com/gambarini/grpcdemo/servers/chat/internal/db"
	"github.com/gambarini/grpcdemo/clients/contact"
)

func main() {

	grpcServer := grpc.NewServer()

	contactClient, conn := contact.NewInternalContactClient()

	defer conn.Close()

	pb.RegisterChatServer(grpcServer, &server.ChatServer{
		DB:            db.NewDB(),
		ContactClient: contactClient,
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 30001))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening on port 30001")
	grpcServer.Serve(listener)

}
