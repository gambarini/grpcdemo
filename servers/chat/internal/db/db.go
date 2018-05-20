package db

import (
	pb "github.com/gambarini/grpcdemo/pb/chat"
	"errors"
)

var (
	streamStore = make(map[string]pb.Chat_StartChatServer)
	ErrNoChatStreamForContact = errors.New("failed to get chat stream for contact")
)

func StoreChatStream(contactID string, chatStream pb.Chat_StartChatServer) {

	_, ok := streamStore[contactID]

	if !ok {
		streamStore[contactID] = chatStream
	}

}

func RemoveChatStream(contactID string) {

	delete(streamStore, contactID)
}

func GetChatStreamByContactID(contactID string) (chatStream pb.Chat_StartChatServer, err error) {

	chatStream, ok := streamStore[contactID]

	if !ok {
		return chatStream, ErrNoChatStreamForContact
	}

	return chatStream, nil
}