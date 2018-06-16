package cliutils

import (
	"google.golang.org/grpc"
	"log"
	"golang.org/x/net/context"
)

const (

	ExternalDomain    = "35.189.51.177:443"

	InternalChatServiceName    = "chat-service.default.svc.cluster.local:50051"
	InternalContactServiceName = "contact-service.default.svc.cluster.local:50051"
	InternalMessageServiceName = "message-service.default.svc.cluster.local:50051"
)

func StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

	cStream, err :=  streamer(ctx, desc, cc, method, opts...)

	log.Printf("Stream %s : %v - %v", method, opts, desc)

	cStream.Trailer()

	return cStream, err
}

func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {


	err := invoker(ctx, method, req, reply, cc, opts...)

	log.Printf("Unary %s : %v - %v -> %v", method, opts, req, reply)


	return err

}
