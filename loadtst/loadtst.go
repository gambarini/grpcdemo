package main

import (
	"golang.org/x/net/context"
	"github.com/gambarini/grpcdemo/cliutils/chat"
	"log"
	"io"
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
	"sync"
)

var (
	MsgPerContactTotal = 500
	ChatTotal          = 100
	MsgPerSec          = time.Second / 3
)

func main() {

	var wg sync.WaitGroup

	for i := 0; i < ChatTotal; i++ {
		wg.Add(1)
		go chatExecute(fmt.Sprintf("A%d", i), fmt.Sprintf("B%d", i), i, &wg)

		time.Sleep(time.Second / 2)
	}

	wg.Wait()
	log.Println("End")
}

func chatExecute(AID, BID string, idx int,  wg *sync.WaitGroup) {

	ctx := context.TODO()

	chatClientA, connA := chat.NewExternalChatClient(true)
	chatClientB, connB := chat.NewExternalChatClient(true)

	defer connA.Close()
	defer connB.Close()
	defer wg.Done()

	sTime := time.Now()

	streamA, err := chatClientA.StartChat(ctx)

	if err != nil {
		log.Printf("[%d] failed to start chat A: %v", idx, err)
		return
	}

	streamB, err := chatClientB.StartChat(ctx)

	if err != nil {
		log.Printf("[%d] failed to start chat B: %v", idx, err)
		return
	}

	log.Printf("[%d] Started chat %s - %s", idx, AID, BID)

	err = connect(AID, BID, streamA, streamB)

	if err != nil {
		log.Printf("[%d] failed to connect: %v", idx, err)
		return
	}

	waitA := make(chan int)
	waitB := make(chan int)

	go Receive(waitA, idx, streamA)
	go Receive(waitB, idx, streamB)

	go Send(AID, BID, idx, streamA)
	go Send(BID, AID, idx, streamB)

	tA := <-waitA
	tB := <-waitB

	eTime := time.Now()

	log.Printf("[%d] End chat %s(%d) - %s(%d) : %s -> %s", idx, AID, tA, BID, tB, sTime.Format("15:04:05"), eTime.Format("15:04:05"))

}

func Send(from, to string, idx int, stream chatpb.Chat_StartChatClient) {

	defer stream.CloseSend()

	for i := 0; i < MsgPerContactTotal; i++ {
		err := stream.Send(&chatpb.Message{
			Text:          fmt.Sprintf("%d", i),
			FromContactId: from,
			ToContactId:   to,
			Type:          chatpb.MessageType_TEXT,
			Timestamp:     &timestamp.Timestamp{Seconds: time.Now().UTC().Unix()},
		})

		if err != nil {
			log.Printf("[%d] failed to send: %v", idx, err)
			return
		}

		time.Sleep(MsgPerSec)
	}
}

func connect(AID, BID string, streamA, streamB chatpb.Chat_StartChatClient) error {
	err := streamA.Send(&chatpb.Message{
		Text:          "",
		FromContactId: AID,
		ToContactId:   BID,
		Type:          chatpb.MessageType_CONNECT,
		Timestamp:     &timestamp.Timestamp{Seconds: time.Now().UTC().Unix()},
	})

	if err != nil {
		return err
	}

	err = streamB.Send(&chatpb.Message{
		Text:          "",
		FromContactId: BID,
		ToContactId:   AID,
		Type:          chatpb.MessageType_CONNECT,
		Timestamp:     &timestamp.Timestamp{Seconds: time.Now().UTC().Unix()},
	})

	if err != nil {
		return err
	}

	return nil
}

func Receive(wait chan int, idx int, stream chatpb.Chat_StartChatClient) {

	c := 0

	for {
		_, err := stream.Recv()

		if err == io.EOF {
			//log.Printf("receive EOF: %v", err)
			wait <- c
			close(wait)
			return
		}

		if err != nil {
			log.Printf("[%d] failed to receive: %v", idx, err)
			wait <- c
			close(wait)
			return
		} else {

			c++

		}

	}
}
