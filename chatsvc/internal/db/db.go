package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"github.com/satori/go.uuid"
)

type DB struct {
	Session *mgo.Session
}

func NewDB(session *mgo.Session) *DB {
	return &DB{session}

}

type ChatConnection struct {
	ID           string `bson:"_id,omitempty"`
	ContactID    string
	Source       string
}

func (db *DB) AddChatConnection(contactID, Source string) (chatConnectionID string, err error) {

	collection := db.Session.DB("grpcdemo").C("Connection")

	id, _ := uuid.NewV4()
	chatConnectionID = fmt.Sprintf("%s-%s", contactID, id)

	connection := &ChatConnection{

		ID: chatConnectionID,
		ContactID: contactID,
		Source: Source,
	}

	err = collection.Insert(connection)

	if err != nil {
		return chatConnectionID, err
	}

	return chatConnectionID, nil
}

func (db *DB) FindContactChatConnectionIDs(contactID string) (connections []ChatConnection, err error) {

	collection := db.Session.DB("grpcdemo").C("Connection")

	err = collection.Find(bson.M{"contactid": contactID}).All(&connections)

	if err != nil {
		return connections, err
	}

	return connections, nil
}

func (db *DB) RemoveChatConnection(chatConnectionID string) (err error) {

	collection := db.Session.DB("grpcdemo").C("Connection")

	err = collection.RemoveId(chatConnectionID)

	if err != nil {
		return err
	}

	return nil
}
