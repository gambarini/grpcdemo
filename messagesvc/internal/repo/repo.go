package repo

import "github.com/gambarini/grpcdemo/dbutils"

type MessageRepository struct {
	DB *dbutils.DB
}

func NewMessageRepository(db *dbutils.DB) *MessageRepository {

	return &MessageRepository{db}

}
