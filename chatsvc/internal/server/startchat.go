package server

import (
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"io"
	"log"
	"fmt"
	"github.com/streadway/amqp"
	"encoding/json"
)

func (server *ChatServer) StartChat(stream chatpb.Chat_StartChatServer) error {

	var streamContactID string
	disconnect := make(chan int)

	for {

		msg, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Received EOF: %s", err)
			disconnect <- 1
			return nil
		}

		if err != nil {
			log.Printf("Error Receiving: %s", err)
			disconnect <- 1
			return err
		}

		log.Printf("Stream Type: %s, From: %s, To: %s, Text: %s", msg.Type, msg.FromContactId, msg.ToContactId, msg.Text)

		streamContactID = msg.FromContactId


		switch msg.Type {

		case chatpb.MessageType_CONNECT:

			err = server.ChatMQ.ReceiveFromQueue(streamContactID, Receive, disconnect, stream)

			if err != nil {
				disconnect <- 1
				return fmt.Errorf("failed to receive for contact ID %s, %s", streamContactID, err)
			}

		case chatpb.MessageType_DISCONNECT:

			disconnect <- 1
			return io.EOF

		case chatpb.MessageType_TEXT:

			toStreamContactID := msg.ToContactId

			err := server.ChatMQ.Send(msg)

			if err != nil {
				log.Printf("Failed to send to contact ID %s, %s", toStreamContactID, err)
			}

		default:
			log.Printf("Unknow message type, %s", msg.Type, )
		}
	}
}

func Receive(deliveries <-chan amqp.Delivery, disconnect <-chan int, stream chatpb.Chat_StartChatServer, channel *amqp.Channel) {

	for {
		select {
			case delivery := <- deliveries:

				var msg chatpb.Message

				err := json.Unmarshal(delivery.Body, &msg)

				if err != nil {
					log.Printf("Error unmarshaling message: %s", err)
					break
				}

				err = stream.Send(&msg)

				if err != nil {
					log.Printf("Error sending message to stream: %s", err)
					break
				}

			case <- disconnect:

				channel.Close()
				return
		}
	}
}

func Cont(){
	/*ctx := context.TODO()

	contactsStream, err := server.ContactClient.StoreContacts(ctx)

	if err != nil {
		log.Printf("Error Storing contact: %s", err)

		return err
	}

	err = contactsStream.Send(&contactpb.Contact{
		Id:   streamContactID,
		Type: contactpb.ContactType_STANDARD,
		Name: "NONE",
	})

	if err != nil {
		log.Printf("Error sending contact: %s", err)

		return err
	}

	contactsStream.CloseSend()*/
}