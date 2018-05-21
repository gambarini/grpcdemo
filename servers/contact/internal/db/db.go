package db

import (
	pb "github.com/gambarini/grpcdemo/pb/contact"
	"errors"
)

var (
	ErrContactNotStored = errors.New("contact is not stored")
)

type DB struct {
	storeContacts map[string]*pb.Contact
}

func NewDB() *DB {
	return &DB{make(map[string]*pb.Contact)}

}

func (db *DB) StoreContact(contact *pb.Contact) {

	db.storeContacts[contact.Id] = contact
}

func (db *DB) FindContact(id string) (contact *pb.Contact, err error) {

	contact, ok := db.storeContacts[id]

	if !ok {
		return contact, ErrContactNotStored
	}

	return contact, nil
}
