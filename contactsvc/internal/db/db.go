package db

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrContactNotStored = errors.New("contact is not stored")
)

type DB struct {
	Session *mgo.Session
}

func NewDB(session *mgo.Session) *DB {

	// Optional. Switch the Session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	return &DB{session}

}

func (db *DB) StoreContact(contact *contactpb.Contact) error {

	storeContacts := db.Session.DB("grpcdemo").C("contact")

	return storeContacts.Insert(contact)

}

func (db *DB) FindContact(id string) (contact *contactpb.Contact, err error) {

	storeContacts := db.Session.DB("grpcdemo").C("contact")

	contact = &contactpb.Contact{}
	err = storeContacts.Find(bson.M{"id": id}).One(contact)

	if err != nil {
		return contact, err
	}

	return contact, nil
}
