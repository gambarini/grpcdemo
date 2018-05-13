package main

import (
	"google.golang.org/grpc"
	"log"
	chat "github.com/gambarini/grpcdemo/chat/pb"
	"golang.org/x/net/context"

	"time"
	"io"
	"flag"
)

var (
	/*Id = flag.Int("ID", 1, "Person ID")
	Name = flag.String("Name", "Mr. One", "Person Name")
	ID = int32(*Id)

	toId = flag.Int("toID", 2, "Person ID")
	toName = flag.String("toName", "Mr. Two", "Person Name")
	toID = int32(*Id)*/

	Id = flag.Int("ID", 2, "Person ID")
	Name = flag.String("Name", "Mr. Two", "Person Name")
	ID = int32(*Id)

	toId = flag.Int("toID", 1, "Person ID")
	toName = flag.String("toName", "Mr. One", "Person Name")
	toID = int32(*Id)
)

func main() {

	flag.Parse()

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:9000", opts...)

	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	defer conn.Close()

	chatClient := chat.NewChatClient(conn)

	ctx := context.TODO()

	stream, err := chatClient.StartChat(ctx)

	if err != nil {
		log.Fatalf("failed to start chat: %v", err)
	}

	wait := make(chan interface{})

	go Receive(wait, stream)

	ConnectToServer(ID, *Name, stream)

	for i := 0; i < 1000; i++ {

		toPerson := &chat.Person{
			Id: toID,
			Name: *toName,
		}

		SendText("from " + *Name, toPerson, stream)

		time.Sleep(time.Second)
	}

	stream.CloseSend()

	<-wait
}

func SendText(text string, toPerson *chat.Person, stream chat.Chat_StartChatClient){

	err := stream.Send(&chat.Message{
		Text: text,
		Person: &chat.Person{
			Id:   ID,
			Name: *Name,
		},
		ToPerson: toPerson,
		Type: chat.Type_TEXT,
	})

	if err != nil {
		log.Fatalf("failed to send: %v", err)
	}
}

func ConnectToServer(ID int32, Name string, stream chat.Chat_StartChatClient) {

	err := stream.Send(&chat.Message{
		Person: &chat.Person{
			Id:   ID,
			Name: Name,
		},
		Type: chat.Type_CONNECT,
	})

	if err != nil {
		log.Fatalf("failed to send: %v", err)
	}
}

func Receive(wait chan interface{}, stream chat.Chat_StartChatClient) {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			close(wait)
			return
		}

		if err != nil {
			log.Fatalf("failed to receive: %v", err)
		}

		log.Printf("Msg: %s, From: %s, To: %s", msg.Text, msg.Person.Name, msg.ToPerson.Name)
	}
}
