package server

import (
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"io"
	"log"
	"fmt"
	"github.com/streadway/amqp"
	"encoding/json"
	"github.com/gambarini/grpcdemo/chatsvc/internal/db"
)

func (server *ChatServer) StartChat(stream chatpb.Chat_StartChatServer) error {

	var streamContactID, chatConnectionID string
	disconnect := make(chan int)

	for {

		msg, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Received EOF: %s", err)
			endChat(disconnect, server.DB, chatConnectionID)
			return nil
		}

		if err != nil {
			log.Printf("Error Receiving: %s", err)
			endChat(disconnect, server.DB, chatConnectionID)
			return err
		}

		log.Printf("Stream Type: %s, From: %s, To: %s, Text: %s", msg.Type, msg.FromContactId, msg.ToContactId, msg.Text)

		streamContactID = msg.FromContactId


		switch msg.Type {

		case chatpb.MessageType_CONNECT:

			chatConnectionID, err = server.DB.AddChatConnection(streamContactID, "CLI-DFL")

			if err != nil {
				endChat(disconnect, server.DB, chatConnectionID)
				return fmt.Errorf("failed to add connection for contact ID %s, %s", streamContactID, err)
			}

			err = server.ChatMQ.ReceiveFromQueue(chatConnectionID, Receive, disconnect, stream)

			if err != nil {
				endChat(disconnect, server.DB, chatConnectionID)
				return fmt.Errorf("failed to receive for contact ID %s, %s", streamContactID, err)
			}

		case chatpb.MessageType_DISCONNECT:

			endChat(disconnect, server.DB, chatConnectionID)
			return io.EOF

		case chatpb.MessageType_TEXT:

			toStreamContactID := msg.ToContactId

			chatConnections, err := server.DB.FindContactChatConnectionIDs(toStreamContactID)

			if err != nil {
				log.Printf("Error finding connections for contact ID %s, %s", toStreamContactID, err)
			}

			if len(chatConnections) == 0 {
				stream.Send(&chatpb.Message{
					Type: chatpb.MessageType_TEXT,
					ToContactId: streamContactID,
					FromContactId: toStreamContactID,
					Text: fmt.Sprintf("Contact %s is not connected. Message cannot be delivered.", toStreamContactID),
				})
			}

			for _, chatConnection := range chatConnections {
				err := server.ChatMQ.Send(chatConnection.ID, msg)

				if err != nil {
					log.Printf("Failed to send to contact ID %s connection %s, %s", chatConnection.ContactID, chatConnection.ID, err)
				}
			}

		default:
			log.Printf("Unknow message type, %s", msg.Type, )
		}
	}
}

func endChat(disconnect chan int, db *db.DB, chatConnectionID string) {

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