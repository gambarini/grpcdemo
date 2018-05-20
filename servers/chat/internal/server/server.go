package server

import (
	pb "github.com/gambarini/grpcdemo/pb/chat"
	"io"
	"github.com/gambarini/grpcdemo/servers/chat/internal/db"
	"log"
)

type ChatServer struct {}

func (server *ChatServer) StartChat(stream pb.Chat_StartChatServer) error {

	var streamContactID string

	for {

		msg, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Error Receiving: %s", err)
			db.RemoveChatStream(streamContactID)
			return nil
		}

		if err != nil {
			log.Printf("Error Receiving: %s", err)
			db.RemoveChatStream(streamContactID)
			return err
		}

		log.Printf("Stream Type: %s, From: %s, To: %s, Text: %s", msg.Type, msg.FromContactId, msg.ToContactId, msg.Text)

		streamContactID = msg.FromContactId

		switch msg.Type {

		case pb.MessageType_CONNECT:

			db.StoreChatStream(streamContactID, stream)

		case pb.MessageType_DISCONNECT:

			db.RemoveChatStream(streamContactID)
			return nil

		case pb.MessageType_TEXT:

			toStreamContactID := msg.ToContactId

			toStream, err := db.GetChatStreamByContactID(toStreamContactID)

			if err != nil {
				log.Printf("Failed to get contact ID %s stream, %s", toStreamContactID, err)
			} else {
				toStream.Send(msg)
			}

		default:
			log.Printf("Unknow message type, %s", msg.Type,)
		}


	}
}
