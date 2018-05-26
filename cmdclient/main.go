package main

import (
	clients "github.com/gambarini/grpcdemo/cliutils/chat"
	"context"
	"fmt"
	"bufio"
	"os"
	"log"
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"io"

	"strings"
	"errors"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

const (
	endKeyword = "/end"
)

var (
	ErrDisconnect = errors.New("disconnecting")
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your ID:")

	ID := readInput(reader)

	fmt.Println("Sent to ID:")

	toID := readInput(reader)

	ctx := context.TODO()

	chatClient, conn := clients.NewExternalChatClient()

	defer conn.Close()

	stream, err := chatClient.StartChat(ctx)

	if err != nil {
		log.Fatalf("failed to start chat: %v", err)
	}

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
		log.Fatalf("Erro connecting: %v", err)
	}

	for {
		fmt.Println("Type '/end' to disconnect")
		fmt.Printf("Type text to %s: \n", toID)

		text := readInput(reader)

		err = Send(stream, text, ID, toID)

		if err == ErrDisconnect {
			return
		}

		if err != nil {
			log.Fatalf("failed to send msg: %v", err)
		}
	}

	stream.CloseSend()

	<-wait

}

func Send(stream chatpb.Chat_StartChatClient, text, ID, toID string) error {

	var err error

	switch text {

	case endKeyword:
		stream.CloseSend()
		return ErrDisconnect

	default:
		err = stream.Send(&chatpb.Message{
			Text:          text,
			FromContactId: ID,
			ToContactId:   toID,
			Type:          chatpb.MessageType_TEXT,
			Timestamp:     &timestamp.Timestamp{Seconds: time.Now().UTC().Unix()},
		})
	}

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

		if msg.Type == chatpb.MessageType_ECHO {
			log.Printf("Sent to %s: %s \n", msg.ToContactId, msg.Text)
		} else {
			log.Printf("Received from %s: %s \n", msg.FromContactId, msg.Text)
		}
	}
}

func readInput(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	return text
}
