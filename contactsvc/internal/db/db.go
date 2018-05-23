package db

import (
	pb "github.com/gambarini/grpcdemo/pb/contact"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


var (
	ErrContactNotStored = errors.New("contactsvc is not stored")
)

type DB struct {
	Session *mgo.Session
}

func NewDB(session *mgo.Session) *DB {

	// Optional. Switch the Session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	return &DB{session}

}

func (db *DB) StoreContact(contact *pb.Contact) error {

	storeContacts := db.Session.DB("grpcdemo").C("contactsvc")

	return storeContacts.Insert(contact)

}

func (db *DB) FindContact(id string) (contact *pb.Contact, err error) {

	storeContacts := db.Session.DB("grpcdemo").C("contactsvc")

	contact = &pb.Contact{}
	err = storeContacts.Find(bson.M{"id": id}).One(contact)

	if err != nil {
		return contact, err
	}

	return contact, nil
}
