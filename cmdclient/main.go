package main

import (
	"github.com/gambarini/grpcdemo/cliutils/chat"
	"github.com/gambarini/grpcdemo/cliutils/message"
	"context"
	"fmt"
	"bufio"
	"os"
	"log"
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"io"

	"strings"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
	"github.com/gambarini/grpcdemo/pb/messagepb"
	"google.golang.org/grpc"
	"os/exec"
	"github.com/gambarini/grpcdemo/cliutils"
)

const (
	endKeyword    = "/end"
	changeKeyword = "/to"
)

var (
	messageClient  messagepb.MessageClient
	ID, toID       string
	ctx            context.Context
	messageConn    *grpc.ClientConn
	reader         *bufio.Reader
	bufferMessages []string

)

func main() {

	grpc.EnableTracing = true

	reader = bufio.NewReader(os.Stdin)

	fmt.Println("Enter your ID:")

	ID = readInput(reader)

	fmt.Println("to ID:")

	toID = readInput(reader)

	ctx = context.Background()

	chatClient, conn := chat.NewExternalChatClient(cliutils.Dial)

	messageClient, messageConn = message.NewExternalMessageClient(cliutils.Dial)

	defer conn.Close()
	defer messageConn.Close()

	stream, err := chatClient.StartChat(ctx)

	if err != nil {
		log.Fatalf("failed to start chat: %v", err)
	}

	initializeMsgBuffer()

	wait := make(chan interface{})

	go Receive(wait, stream)

	err = stream.Send(&chatpb.Message{
		Text:          "",
		FromContactId: ID,
		ToContactId:   toID,
		Type:          chatpb.MessageType_CONNECT,
		Timestamp:     &timestamp.Timestamp{Seconds: time.Now().UTC().Unix()},
	})

	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}

	clear()
	display()
	desc()

	for {

		text := askCommand()

		switch text {

		case endKeyword:
			stream.CloseSend()
			return

		case changeKeyword:
			fmt.Println("to ID:")
			toID = readInput(reader)
			initializeMsgBuffer()
			clear()
			display()
			desc()

		default:
			err = Send(stream, text)
		}

		if err != nil {
			log.Fatalf("failed to send msg: %v", err)
		}

	}

	<-wait

}

func desc() {
	fmt.Println("------------------------------")
	fmt.Println("Type '/end' to disconnect")
	fmt.Println("Type '/to' to chat to a new contact")
	fmt.Printf("%s type text to %s: \n", ID, toID)

}

func askCommand() string {

	return readInput(reader)
}

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func display() {
	for _, m := range bufferMessages {
		fmt.Println(m)
	}
}

func initializeMsgBuffer() {

	bufferMessages = make([]string, 0)

	stream, _ := messageClient.GetMessages(ctx, &messagepb.Filter{
		ContactId: ID,
		FromTimestamp: &timestamp.Timestamp{
			Seconds: 0,
		},
	})

	for {
		msg, err := stream.Recv()

		if err != nil {

			break
		}

		if msg.FromContactId == toID || msg.ToContactId == toID {

			updateMsgBuffer(msg)
		}
	}

}

func updateMsgBuffer(msg *chatpb.Message) {
	bufferMessages = append(bufferMessages, fmt.Sprintf("%s: %s \n", msg.FromContactId, msg.Text))
}

func Send(stream chatpb.Chat_StartChatClient, text string) error {

	var err error

	err = stream.Send(&chatpb.Message{
		Text:          text,
		FromContactId: ID,
		ToContactId:   toID,
		Type:          chatpb.MessageType_TEXT,
		Timestamp:     &timestamp.Timestamp{Seconds: time.Now().UTC().Unix()},
	})

	if err != nil {
		return err
	}

	return nil
}

func Receive(wait chan interface{}, stream chatpb.Chat_StartChatClient) {

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			close(wait)
			return
		}

		if err != nil {
			log.Fatalf("failed to receive: %v", err)
		}

		if msg.Type == chatpb.MessageType_SYSTEM {
			fmt.Printf(" *** %s ***\n", msg.Text)
		}

		if msg.Type == chatpb.MessageType_TEXT && msg.FromContactId == toID {
			updateMsgBuffer(msg)
			clear()
			display()
			desc()
		}

		if msg.Type == chatpb.MessageType_ECHO && msg.ToContactId == toID {
			updateMsgBuffer(msg)
			clear()
			display()
			desc()
		}

	}
}

func readInput(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	return text
}
