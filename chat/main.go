package main

import (
	"github.com/gambarini/grpcdemo/chat/pb"
	"google.golang.org/grpc"
	"fmt"
	"log"
	"net"
	"io"
)

var (
	streamStore map[int32]pb.Chat_StartChatServer
)

type chatServer struct {}

func (server *chatServer) StartChat(stream pb.Chat_StartChatServer) error {

	var streamID int32

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			clearStreamStore(streamID)
			return nil
		}

		if err != nil {
			clearStreamStore(streamID)
			return err
		}

		log.Printf("Msg: %s, From %d", msg.Type, msg.Person.Id)

		switch msg.Type {
		case pb.Type_CONNECT:
			streamID = msg.Person.Id
			resolveStreamStore(streamID, stream)

		case pb.Type_DISCONNECT:
			clearStreamStore(streamID)
			return nil

		case pb.Type_TEXT:
			toID := msg.ToPerson.Id
			toStream := resolveStreamStore(toID, nil)

			if toStream != nil {
				err = toStream.Send(msg)

				if err != nil {
					return err
				}
			}
		default:

		}
	}
}

func resolveStreamStore(ID int32, stream pb.Chat_StartChatServer) (storedStream pb.Chat_StartChatServer) {

	storedStream, ok := streamStore[ID]

	if !ok || stream != nil {
		streamStore[ID] = stream
		storedStream = stream
	}

	return storedStream
}

func clearStreamStore(ID int32) {

	delete(streamStore, ID)
}

func main() {

	//var opts []grpc.ServerOption

	streamStore = make(map[int32]pb.Chat_StartChatServer)

	grpcServer := grpc.NewServer()

	pb.RegisterChatServer(grpcServer, &chatServer{})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening on :9000")
	grpcServer.Serve(listener)
}
