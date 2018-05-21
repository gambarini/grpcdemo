package server

import (
	pb "github.com/gambarini/grpcdemo/pb/chat"
	contactPb "github.com/gambarini/grpcdemo/pb/contact"
	"io"
	"github.com/gambarini/grpcdemo/servers/chat/internal/db"
	"log"
	"golang.org/x/net/context"
)

type ChatServer struct {
	DB            *db.DB
	ContactClient contactPb.ContactsClient
}

func (server *ChatServer) StartChat(stream pb.Chat_StartChatServer) error {

	var streamContactID string

	for {

		msg, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Received EOF: %s", err)
			server.DB.RemoveChatStream(streamContactID)
			return nil
		}

		if err != nil {
			log.Printf("Error Receiving: %s", err)
			server.DB.RemoveChatStream(streamContactID)
			return err
		}

		log.Printf("Stream Type: %s, From: %s, To: %s, Text: %s", msg.Type, msg.FromContactId, msg.ToContactId, msg.Text)

		streamContactID = msg.FromContactId

		ctx := context.TODO()

		contactsStream, err := server.ContactClient.StoreContacts(ctx)

		if err != nil {
			log.Printf("Error Storing contact: %s", err)
			server.DB.RemoveChatStream(streamContactID)
			return err
		}

		err = contactsStream.Send(&contactPb.Contact{
			Id:   streamContactID,
			Type: contactPb.ContactType_STANDARD,
			Name: "NONE",
		})

		if err != nil {
			log.Printf("Error sending contact: %s", err)
			server.DB.RemoveChatStream(streamContactID)
			return err
		}

		contactsStream.CloseSend()

		switch msg.Type {

		case pb.MessageType_CONNECT:

			server.DB.StoreChatStream(streamContactID, stream)

		case pb.MessageType_DISCONNECT:

			server.DB.RemoveChatStream(streamContactID)
			return io.EOF

		case pb.MessageType_TEXT:

			toStreamContactID := msg.ToContactId

			toStream, err := server.DB.GetChatStreamByContactID(toStreamContactID)

			if err != nil {
				log.Printf("Failed to get contact ID %s stream, %s", toStreamContactID, err)
			} else {
				toStream.Send(msg)
			}

		default:
			log.Printf("Unknow message type, %s", msg.Type, )
		}

	}
}
