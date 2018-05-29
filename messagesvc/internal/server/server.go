package server

import (
	"github.com/gambarini/grpcdemo/pb/messagepb"
	"github.com/gambarini/grpcdemo/messagesvc/internal/repo"
	"io"
	"log"
	"fmt"
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type MessageServer struct {
	Repository *repo.MessageRepository
}

func (server *MessageServer) StoreMessages(stream messagepb.Message_StoreMessagesServer) error {

	for {

		storeMessage, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Received EOF: %s", err)
			return nil
		}

		if err != nil {
			log.Printf("Error Receiving: %s", err)
			return fmt.Errorf("error receiving: %s", err)
		}

		_, err = server.Repository.Store(storeMessage.ContactId, repo.Message{
			Seconds:       storeMessage.Message.Timestamp.Seconds,
			Type:          storeMessage.Message.Type,
			Text:          storeMessage.Message.Text,
			FromContactID: storeMessage.Message.FromContactId,
			ToContactID:   storeMessage.Message.ToContactId,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (server *MessageServer) GetMessages(filter *messagepb.Filter, stream messagepb.Message_GetMessagesServer) (err error) {

	items := make(chan repo.Item, 10)
	abort := make(chan bool)

	defer close(abort)

	go server.Repository.GetMessages(filter.ContactId, filter.FromTimestamp.Seconds, items, abort)

	for item := range items {

		if item.Err != nil {
			return item.Err
		}

		err := stream.Send(&chatpb.Message{
			ToContactId:   item.Msg.ToContactID,
			FromContactId: item.Msg.FromContactID,
			Text:          item.Msg.Text,
			Type:          item.Msg.Type,
			Timestamp:     &timestamp.Timestamp{Seconds: item.Msg.Seconds},
		})

		if err != nil {
			abort <- true
			return err
		}
	}

	return nil

}
