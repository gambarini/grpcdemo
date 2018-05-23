package main

import (
	clients "github.com/gambarini/grpcdemo/clients/chat"
	"context"
	"fmt"
	"bufio"
	"os"
	"log"
	"github.com/gambarini/grpcdemo/pb/chat"
	"io"

	"strings"
	"errors"
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
		log.Fatalf("failed to start chatsvc: %v", err)
	}

	wait := make(chan interface{})

	go Receive(wait, stream)

	err = stream.Send(&chat.Message{
		Text:          "",
		FromContactId: ID,
		ToContactId:   toID,
		Type:          chat.MessageType_CONNECT,
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

func Send(stream chat.Chat_StartChatClient, text, ID, toID string) error {

	var err error

	switch text {

	case endKeyword:
		stream.CloseSend()
		return ErrDisconnect

	default:
		err = stream.Send(&chat.Message{
			Text:          text,
			FromContactId: ID,
			ToContactId:   toID,
			Type:          chat.MessageType_TEXT,
		})
	}

	if err != nil {
		return err
	}

	return nil
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

		log.Printf("Received from %s: %s \n", msg.FromContactId, msg.Text)
	}
}

func readInput(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	return text
}
