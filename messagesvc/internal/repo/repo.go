package repo

import (
	"github.com/gambarini/grpcdemo/dbutils"

	"github.com/gambarini/grpcdemo/pb/chatpb"

	"gopkg.in/mgo.v2"
	"github.com/satori/go.uuid"

	"gopkg.in/mgo.v2/bson"
	"log"
)

type (
	MessageRepository struct {
		DB dbutils.DB
	}

	Message struct {
		ID            string `bson:"_id,omitempty"`
		ToContactID   string
		FromContactID string
		Seconds       int64
		Text          string
		Type          chatpb.MessageType
	}

	Item struct {
		Msg Message
		Err error
	}
)

func (repo *MessageRepository) GetCollection(contactID string) (*mgo.Session, *mgo.Collection) {

	session := repo.DB.GetSession()

	return session, session.DB("message").C(contactID)
}

func NewMessageRepository(db dbutils.DB) *MessageRepository {

	return &MessageRepository{db}
}

func (repo *MessageRepository) Store(contactID string, message Message) (Message, error) {

	session, collection := repo.GetCollection(contactID)

	defer session.Close()

	u2, err := uuid.NewV4()

	if err != nil {
		return message, err
	}

	message.ID = u2.String()

	err = collection.Insert(message)

	if err != nil {
		return message, err
	}

	return message, nil

}

func (repo *MessageRepository) GetMessages(contactID string, lastSeconds int64, items chan Item, abort chan bool) {

	session, collection := repo.GetCollection(contactID)

	defer session.Close()

	iter := collection.Find(bson.M{"seconds": bson.M{"$gt": lastSeconds}}).Sort("seconds").Iter()

	defer iter.Close()
	defer close(items)

	go func() {
		end := <-abort

		if end {
			log.Printf("Aborting iterator on get messages")
			iter.Close()
		}
	}()

	for {

		err := iter.Err()

		if err != nil {
			log.Printf("Erro on iterator get messages, %s", err)
			items <- Item{
				Err: err,
			}
			break
		}

		if iter.Done() {
			log.Printf("Done on iterator get messages")
			break
		}

		var msg Message
		iter.Next(&msg)

		items <- Item{
			Msg: msg,
			Err: nil,
		}
	}

	log.Printf("Ending get messages")

}
