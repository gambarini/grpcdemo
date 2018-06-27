package server

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"github.com/gambarini/grpcdemo/contactsvc/internal/repo"
	"io"
	"log"
	"github.com/gambarini/grpcdemo/dbutils"
	"github.com/gambarini/grpcdemo/svcutils"
)

type ContactsServer struct {
	ContactRepository *repo.ContactRepository
}

func (server *ContactsServer) Initialize(main *svcutils.Main) (err error) {

	db, err := dbutils.NewMongoDB(dbutils.MongoDBURL)

	if err != nil {
		return err
	}

	server.ContactRepository = repo.NewContactRepository(db)

	contactpb.RegisterContactsServer(main.GRPCServer, server)

	return nil
}

func (server *ContactsServer) CleanUp() {

	server.ContactRepository.DB.CleanUp()
}

func (server *ContactsServer) StoreContacts(stream contactpb.Contacts_StoreContactsServer) error {

	for {

		store, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Received EOF: %s", err)
			return nil
		}

		if err != nil {
			log.Printf("Error Receiving: %s", err)
			return err
		}

		log.Printf("Contact to store: %v", store.Contact)

		return server.ContactRepository.StoreContact(store.ListContactId, repo.Contact{
			ID:   store.Contact.Id,
			Name: store.Contact.Name,
		})

	}
}

func (server *ContactsServer) ListContacts(filterContact *contactpb.Filter, stream contactpb.Contacts_ListContactsServer) error {

	items := make(chan repo.Item, 10)
	abort := make(chan bool)

	defer close(abort)

	go server.ContactRepository.FindContact(filterContact.ListContactId, items, abort)

	for item := range items {

		if item.Err != nil {
			return item.Err
		}

		err := stream.Send(&contactpb.Contact{
			Id:   item.Contact.ID,
			Name: item.Contact.Name,
		})

		if err != nil {
			abort <- true
			return err
		}
	}

	return nil
}
