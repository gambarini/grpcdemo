package server

import (
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"io"
	"github.com/gambarini/grpcdemo/chatsvc/internal/db"
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ChatServer struct {
	DB                *db.DB
	ContactClient     contactpb.ContactsClient
	ContactClientConn *grpc.ClientConn
}

func (server *ChatServer) StartChat(stream chatpb.Chat_StartChatServer) error {

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
			log.Printf("Error Storing contactsvc: %s", err)
			server.DB.RemoveChatStream(streamContactID)
			return err
		}

		err = contactsStream.Send(&contactpb.Contact{
			Id:   streamContactID,
			Type: contactpb.ContactType_STANDARD,
			Name: "NONE",
		})

		if err != nil {
			log.Printf("Error sending contactsvc: %s", err)
			server.DB.RemoveChatStream(streamContactID)
			return err
		}

		contactsStream.CloseSend()

		switch msg.Type {

		case chatpb.MessageType_CONNECT:

			server.DB.StoreChatStream(streamContactID, stream)

		case chatpb.MessageType_DISCONNECT:

			server.DB.RemoveChatStream(streamContactID)
			return io.EOF

		case chatpb.MessageType_TEXT:

			toStreamContactID := msg.ToContactId

			toStream, err := server.DB.GetChatStreamByContactID(toStreamContactID)

			if err != nil {
				log.Printf("Failed to get contactsvc ID %s stream, %s", toStreamContactID, err)
			} else {
				toStream.Send(msg)
			}

		default:
			log.Printf("Unknow message type, %s", msg.Type, )
		}

	}
}
