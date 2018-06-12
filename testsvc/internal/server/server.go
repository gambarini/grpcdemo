package server

import (
	"github.com/gambarini/grpcdemo/pb/testpb"
	"golang.org/x/net/context"
	"github.com/gambarini/grpcdemo/svcutils"
	"io"
	"log"
	"fmt"
	"google.golang.org/grpc/reflection"
)

type (
	TestServer struct{}
)

func (server *TestServer) TestStream(stream testpb.Test_TestStreamServer) error {

	for {

		msg, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Received EOF: %s", err)

			return nil
		}

		if err != nil {
			log.Printf("Error Receiving: %s", err)

			return fmt.Errorf("error receiving: %s", err)
		}

		log.Printf("Received Stream: %v", msg)

		err = stream.Send(&testpb.TestMsg{
			Text: "Reply to " + msg.Text,
		})

		if err != nil {
			log.Printf("Error Sending: %s", err)

		}

	}
}

func (server *TestServer) TestUnary(ctx context.Context, msg *testpb.TestMsg) (*testpb.TestMsg, error) {

	log.Printf("Received Unary: %v", msg)

	return &testpb.TestMsg{
		Text: "Reply to " + msg.Text,
	}, nil
}



func (server *TestServer) Initialize(main *svcutils.Main) error {

	testpb.RegisterTestServer(main.GRPCServer, server)

	reflection.Register(main.GRPCServer)

	return nil
}

func (server *TestServer) CleanUp() {

}