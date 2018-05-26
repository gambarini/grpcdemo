package repo

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/gambarini/grpcdemo/dbutils"
	"gopkg.in/mgo.v2"
)

type ChatRepository struct {
	DB *dbutils.DB
}

func NewChatRepository(db *dbutils.DB) *ChatRepository {

	return &ChatRepository{db}

}

type ChatConnection struct {
	ID         string `bson:"_id,omitempty"`
	ContactID  string
	Source     string
	ConnNumber int
}

func (repo *ChatRepository) GetCollection() *mgo.Collection {

	return repo.DB.Session.DB("grpcdemo").C("Connection")
}

func (repo *ChatRepository) AddChatConnection(contactID, Source string) (chatConnection *ChatConnection, err error) {

	collection := repo.GetCollection()

	chatConnection, err = repo.FindContactChatConnection(contactID)

	if err != nil {
		return chatConnection, err
	}

	if chatConnection == nil {

		id, _ := uuid.NewV4()
		chatConnectionID := fmt.Sprintf("%s-%s", contactID, id)

		chatConnection = &ChatConnection{
			ID:         chatConnectionID,
			ContactID:  contactID,
			Source:     Source,
			ConnNumber: 1,
		}

	} else {

		chatConnection.ConnNumber += 1

	}

	_, err = collection.UpsertId(chatConnection.ID, chatConnection)

	if err != nil {
		return chatConnection, err
	}

	return chatConnection, nil
}

func (repo *ChatRepository) FindContactChatConnection(contactID string) (chatConnection *ChatConnection, err error) {

	collection := repo.GetCollection()

	var connections []ChatConnection

	err = collection.Find(bson.M{"contactid": contactID}).All(&connections)

	if err != nil {
		return chatConnection, err
	}

	if len(connections) == 0 {
		return chatConnection, nil
	}

	chatConnection = &connections[0]

	return chatConnection, nil
}

func (repo *ChatRepository) RemoveChatConnection(contactID string) (err error) {

	collection := repo.GetCollection()

	chatConnection, err := repo.FindContactChatConnection(contactID)

	if err != nil {
		return err
	}

	if chatConnection != nil {
		chatConnection.ConnNumber -= 1

		_, err = collection.UpsertId(chatConnection.ID, &chatConnection)

		if err != nil {
			return err
		}
	}

	return nil
}
