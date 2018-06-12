package main

import (
	"google.golang.org/grpc"
	"github.com/gambarini/grpcdemo/cliutils"
	"log"

	"github.com/gambarini/grpcdemo/pb/testpb"
	"context"
	"io"
	"google.golang.org/grpc/credentials"
	"fmt"
	"os"
	"crypto/tls"
)

var (
	authorityCertFile = fmt.Sprintf("%s/src/github.com/gambarini/grpcdemo/certificate/ca.crt", os.Getenv("GOPATH"))
	serverCertFile = fmt.Sprintf("%s/src/github.com/gambarini/grpcdemo/certificate/server.crt", os.Getenv("GOPATH"))
)

func main() {

	creds := credentials.NewTLS(&tls.Config{ InsecureSkipVerify: true})

	var opts []grpc.DialOption

	//opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithStreamInterceptor(cliutils.StreamClientInterceptor))
	opts = append(opts, grpc.WithUnaryInterceptor(cliutils.UnaryClientInterceptor))

	conn, err := grpc.Dial("35.189.55.87:443", opts...)

	if err != nil {
		log.Fatalf("Failed to dial Chat Service: %v", err)
	}

	client := testpb.NewTestClient(conn)

	_, err = client.TestUnary(context.TODO(), &testpb.TestMsg{
		Text: "OK!",
	})

	if err != nil {
		log.Fatalf("Error sending unary: %s", err)
	}

	stream, err := client.TestStream(context.TODO())

	if err != nil {
		log.Fatalf("Erro connecting stream: %s", err)
	}

	wait := make(chan int)

	go func() {
		for {

			rec, err := stream.Recv()

			if err == io.EOF {
				close(wait)
				return
			}

			if err != nil {
				log.Fatalf("failed to receive: %v", err)
			}

			log.Printf("Stream rec: %v", rec)
		}
	}()

	err = stream.Send(&testpb.TestMsg{
		Text: "OK!",
	})

	if err != nil {
		log.Fatalf("Error sending stream: %s", err)
	}

	stream.CloseSend()
	<-wait

}
