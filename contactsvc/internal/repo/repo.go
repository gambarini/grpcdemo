package repo

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"github.com/gambarini/grpcdemo/dbutils"
)

var (
	ErrContactNotStored = errors.New("contact is not stored")
)

type ContactRepository struct {
	DB dbutils.DB
}

func NewContactRepository(db dbutils.DB) *ContactRepository {

	return &ContactRepository{db}

}

func (repo *ContactRepository) StoreContact(contact *contactpb.Contact) error {

	session := repo.DB.GetSession()

	defer session.Close()

	storeContacts := session.DB("contact").C("contact")

	return storeContacts.Insert(contact)

}

func (repo *ContactRepository) FindContact(id string) (contact *contactpb.Contact, err error) {

	session := repo.DB.GetSession()

	defer session.Close()

	storeContacts := session.DB("contact").C("contact")

	contact = &contactpb.Contact{}
	err = storeContacts.Find(bson.M{"id": id}).One(contact)

	if err != nil {
		return contact, err
	}

	return contact, nil
}
