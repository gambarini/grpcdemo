package svcutils

import (
	"net"
	"fmt"
	"log"
	"google.golang.org/grpc"
	"os/signal"
	"os"
	"syscall"
)

type InitializationFunc func(mainServer *MainServer) (err error)
type CleanUpFunc func(mainServer *MainServer)

type MainServer struct {
	Name           string
	ServerPort     int
	Initialization InitializationFunc
	CleanUp        CleanUpFunc
	GRPCServer     *grpc.Server
	ServerObjects  interface{}
	signalChannel  chan os.Signal
}

func (mainServer *MainServer) Run() {

	mainServer.signalChannel = make(chan os.Signal, 1)

	signal.Notify(mainServer.signalChannel, syscall.SIGINT)  // Handling Ctrl + C
	signal.Notify(mainServer.signalChannel, syscall.SIGTERM) // Handling Docker stop

	mainServer.GRPCServer = grpc.NewServer()

	log.Print("Initializing server resources...")
	err := mainServer.Initialization(mainServer)

	if err != nil {
		log.Fatalf("Failed to initialize server resources: %s", err)
	}

	log.Print("Initialization Done!")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", mainServer.ServerPort))

	if err != nil {
		mainServer.CleanUp(mainServer)
		log.Fatalf("Error creating TCP listener on port %d: %s", mainServer.ServerPort, err)
	}

	go mainServer.handleSystemSignals()

	log.Printf("%s listening on port %d", mainServer.Name, mainServer.ServerPort)
	err = mainServer.GRPCServer.Serve(listener)

	if err != nil {
		mainServer.CleanUp(mainServer)
		log.Fatalf("Error while starting to serve: %s", err)
	}

	log.Print("Cleanning up server resorces...")
	mainServer.CleanUp(mainServer)

	log.Print("Server stopped.")
}

func (mainServer *MainServer) handleSystemSignals() {

	sig := <-mainServer.signalChannel

	log.Printf("System signal received: %s", sig.String())

	log.Print("Gracefully stopping server...")
	mainServer.GRPCServer.GracefulStop()

}
