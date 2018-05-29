package server

import (
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"io"
	"log"
	"fmt"
	"github.com/streadway/amqp"
	"encoding/json"
	"github.com/gambarini/grpcdemo/chatsvc/internal/repo"
	"golang.org/x/net/context"
	"github.com/gambarini/grpcdemo/pb/messagepb"
)

func (server *ChatServer) StartChat(stream chatpb.Chat_StartChatServer) error {

	var streamContactID string
	disconnect := make(chan int)

	storeMessageStream, err := server.MessageClient.StoreMessages(context.TODO())

	if err != nil {
		log.Printf("Failed to connect to message store: %s", err)
		endChat(disconnect, server.Repository, streamContactID)
		return fmt.Errorf("failed to connect to message store: %s", err)
	}

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

			deliveries, channel, err := server.ChatMQ.ReceiveFromQueue(chatConnection.ContactID)

			if err != nil {
				log.Printf("failed to start receiving for contact ID %s, %s", streamContactID, err)
				endChat(disconnect, server.Repository, chatConnection.ContactID)
				return fmt.Errorf("failed to start receiving for contact ID %s, %s", streamContactID, err)
			}

			go Receive(deliveries, disconnect, stream, channel)

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

			if chatConnections == nil {
				stream.Send(&chatpb.Message{
					Type: chatpb.MessageType_SYSTEM,
					ToContactId: streamContactID,
					FromContactId: toStreamContactID,
					Text: fmt.Sprintf("Contact %s does not exist. Message cannot be delivered.", toStreamContactID),
				})
				continue
			}

			err = storeMessages(storeMessageStream, msg, stream)

			if err != nil {
				log.Printf("%s", err)

				stream.Send(&chatpb.Message{
					Type: chatpb.MessageType_SYSTEM,
					ToContactId: msg.FromContactId,
					FromContactId: msg.FromContactId,
					Text: fmt.Sprintf("Error: Message could not be delivered due to an internal error."),
				})

				continue
			}

			if chatConnections.ConnNumber > 0 {

				err = server.ChatMQ.Send(chatConnections.ContactID, msg)

				if err != nil {
					log.Printf("Failed to send to contact ID %s, %s", chatConnections.ContactID, err)
					continue
				}

			} else {
				log.Printf("No connections for contact ID %s, The message will be stored only.", chatConnections.ContactID)
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

func storeMessages(storeMessageStream messagepb.Message_StoreMessagesClient, msg *chatpb.Message, stream chatpb.Chat_StartChatServer) (err error) {

	err = storeMessageStream.Send(&messagepb.StoreMessage{
		Message: msg,
		ContactId: msg.FromContactId,
	})

	if err != nil {
		return fmt.Errorf("failed to send message to store for contact ID %s, %s", msg.FromContactId, err)
	}

	err = storeMessageStream.Send(&messagepb.StoreMessage{
		Message: msg,
		ContactId: msg.ToContactId,
	})

	if err != nil {
		return fmt.Errorf("failed to send message to store for contact ID %s, %s", msg.ToContactId, err)
	}

	return nil

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