package svcutils

import (
	"log"
	"os"
	"google.golang.org/grpc"
	"os/signal"
	"syscall"
	"net"
	"fmt"
)

type (
	Server interface {
		Initialize(main *Main) error
		CleanUp()
	}

	Main struct {
		Name          string
		ServerPort    int
		GRPCServer    *grpc.Server
		signalChannel chan os.Signal
		Server        Server
	}
)

func (main *Main) Run() {

	main.signalChannel = make(chan os.Signal, 1)

	signal.Notify(main.signalChannel, syscall.SIGINT)  // Handling Ctrl + C
	signal.Notify(main.signalChannel, syscall.SIGTERM) // Handling Docker stop

	main.GRPCServer = grpc.NewServer()

	log.Print("Initializing main resources...")
	err := main.Server.Initialize(main)

	if err != nil {
		log.Fatalf("Failed to initialize main resources: %s", err)
	}

	log.Print("Initialization Done!")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", main.ServerPort))

	if err != nil {
		main.Server.CleanUp()
		log.Fatalf("Error creating TCP listener on port %d: %s", main.ServerPort, err)
	}

	go main.handleSystemSignals()

	log.Printf("%s listening on port %d", main.Name, main.ServerPort)
	err = main.GRPCServer.Serve(listener)

	if err != nil {
		main.Server.CleanUp()
		log.Fatalf("Error while starting to serve: %s", err)
	}

	log.Print("Cleanning up main resorces...")
	main.Server.CleanUp()

	log.Print("Server stopped.")
}



func (main *Main) handleSystemSignals() {

	sig := <-main.signalChannel

	log.Printf("System signal received: %s", sig.String())

	log.Print("Gracefully stopping server...")
	main.GRPCServer.GracefulStop()

}
