package server

import (
	"github.com/gambarini/grpcdemo/pb/messagepb"
	"github.com/gambarini/grpcdemo/messagesvc/internal/repo"
)

type MessageServer struct {
	Repository        *repo.MessageRepository
}

func (server *MessageServer) SaveMessages(messagepb.Message_SaveMessagesServer) error {

	return nil
}


func (server *MessageServer) GetMessages(*messagepb.Filter, messagepb.Message_GetMessagesServer) error {

	return nil
}