package db

import (
	"github.com/gambarini/grpcdemo/pb/chatpb"
	"errors"
)

var (
	ErrNoChatStreamForContact = errors.New("failed to get chat stream for contact")
)

type DB struct {
	streamStore map[string]chatpb.Chat_StartChatServer
}

func NewDB() *DB {
	return &DB{make(map[string]chatpb.Chat_StartChatServer)}

}

func (db *DB) StoreChatStream(contactID string, chatStream chatpb.Chat_StartChatServer) {

	_, ok := db.streamStore[contactID]

	if !ok {
		db.streamStore[contactID] = chatStream
	}

}

func (db *DB) RemoveChatStream(contactID string) {

	delete(db.streamStore, contactID)
}

func (db *DB) GetChatStreamByContactID(contactID string) (chatStream chatpb.Chat_StartChatServer, err error) {

	chatStream, ok := db.streamStore[contactID]

	if !ok {
		return chatStream, ErrNoChatStreamForContact
	}

	return chatStream, nil
}
