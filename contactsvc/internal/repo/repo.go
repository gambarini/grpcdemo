package repo

import (
	"errors"
	"github.com/gambarini/grpcdemo/dbutils"
	"log"
)

var (
	ErrContactNotStored = errors.New("contact is not stored")
)

type (
	ContactRepository struct {
		DB dbutils.DB
	}

	Contact struct {
		ID   string
		Name string
	}

	Item struct {
		Contact Contact
		Err     error
	}
)

func NewContactRepository(db dbutils.DB) *ContactRepository {

	return &ContactRepository{db}

}

func (repo *ContactRepository) StoreContact(listContactID string, contact Contact) error {

	session := repo.DB.GetSession()

	defer session.Close()

	storeContacts := session.DB("contact").C(listContactID)

	return storeContacts.Insert(contact)

}

func (repo *ContactRepository) FindContact(listContactID string, items chan Item, abort chan bool) {

	session := repo.DB.GetSession()

	defer session.Close()

	storeContacts := session.DB("contact").C(listContactID)

	iter := storeContacts.Find(nil).Iter()

	defer iter.Close()
	defer close(items)

	go func() {
		end := <-abort

		if end {
			log.Printf("Aborting iterator on find contacts")
			iter.Close()
		}
	}()

	for {

		err := iter.Err()

		if err != nil {
			log.Printf("Erro on iterator find contacts, %s", err)
			items <- Item{
				Err: err,
			}
			break
		}

		if iter.Done() {
			log.Printf("Done on iterator find contacts")
			break
		}

		var contact Contact
		iter.Next(&contact)

		items <- Item{
			Contact: contact,
			Err:     nil,
		}

	}

	log.Printf("Ending find contacts")

}
