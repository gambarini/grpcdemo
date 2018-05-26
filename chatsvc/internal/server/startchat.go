package server

import (
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"io"
	"log"
	"fmt"
	"github.com/streadway/amqp"
	"encoding/json"
	"github.com/gambarini/grpcdemo/chatsvc/internal/repo"
)

func (server *ChatServer) StartChat(stream chatpb.Chat_StartChatServer) error {

	var streamContactID string
	disconnect := make(chan int)

	for {

		msg, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Received EOF: %s", err)
			endChat(disconnect, server.Repository, streamContactID)
			return nil
		}

		if err != nil {
			log.Printf("Error Receiving: %s", err)
			endChat(disconnect, server.Repository, streamContactID)
			return fmt.Errorf("error receiving: %s", err)
		}

		log.Printf("Stream Type: %s, From: %s, To: %s, Text: %s", msg.Type, msg.FromContactId, msg.ToContactId, msg.Text)

		streamContactID = msg.FromContactId


		switch msg.Type {

		case chatpb.MessageType_CONNECT:

			chatConnection, err := server.Repository.AddChatConnection(streamContactID, "CLI-DFL")

			if err != nil {
				log.Printf("failed to add connection for contact ID %s, %s", streamContactID, err)
				endChat(disconnect, server.Repository, chatConnection.ContactID)
				return fmt.Errorf("failed to add connection for contact ID %s, %s", streamContactID, err)
			}

			err = server.ChatMQ.ReceiveFromQueue(chatConnection.ContactID, Receive, disconnect, stream)

			if err != nil {
				log.Printf("failed to receive for contact ID %s, %s", streamContactID, err)
				endChat(disconnect, server.Repository, chatConnection.ContactID)
				return fmt.Errorf("failed to receive for contact ID %s, %s", streamContactID, err)
			}

		case chatpb.MessageType_DISCONNECT:

			endChat(disconnect, server.Repository, streamContactID)
			return io.EOF

		case chatpb.MessageType_TEXT:

			toStreamContactID := msg.ToContactId

			chatConnections, err := server.Repository.FindContactChatConnection(toStreamContactID)

			if err != nil {
				log.Printf("Error finding connections for contact ID %s, %s", toStreamContactID, err)
				continue
			}

			if chatConnections == nil || chatConnections.ConnNumber == 0 {
				stream.Send(&chatpb.Message{
					Type: chatpb.MessageType_TEXT,
					ToContactId: streamContactID,
					FromContactId: toStreamContactID,
					Text: fmt.Sprintf("Contact %s is not connected. Message cannot be delivered.", toStreamContactID),
				})
				continue
			}

			err = server.ChatMQ.Send(chatConnections.ContactID, msg)

			if err != nil {
				log.Printf("Failed to send to contact ID %s, %s", chatConnections.ContactID, err)
				continue
			}

			msg.Type = chatpb.MessageType_ECHO

			err = server.ChatMQ.Send(streamContactID, msg)

			if err != nil {
				log.Printf("Failed to send echo message %s, %s", chatConnections.ContactID, err)
				continue
			}

		default:
			log.Printf("Unknow message type, %s", msg.Type)
		}
	}
}

func endChat(disconnect chan int, db *repo.ChatRepository, chatConnectionID string) {

	close(disconnect)

	db.RemoveChatConnection(chatConnectionID)
}

func Receive(deliveries <-chan amqp.Delivery, disconnect <-chan int, stream chatpb.Chat_StartChatServer, channel *amqp.Channel) {

	for {
		select {
			case delivery := <- deliveries:

				var msg chatpb.Message

				err := json.Unmarshal(delivery.Body, &msg)

				if err != nil {
					log.Printf("Error unmarshaling message: %s", err)
					continue
				}

				err = stream.Send(&msg)

				if err != nil {
					log.Printf("Error sending message to stream: %s", err)
					continue
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