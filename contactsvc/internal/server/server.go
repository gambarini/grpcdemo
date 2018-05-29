package server

import (
	"github.com/gambarini/grpcdemo/pb/contactpb"
	"github.com/gambarini/grpcdemo/contactsvc/internal/repo"
	"io"
	"log"
)

type ContactsServer struct {
	ContactRepository *repo.ContactRepository
}

func (server *ContactsServer) StoreContacts(stream contactpb.Contacts_StoreContactsServer) error {

	for {

		contact, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Received EOF: %s", err)
			return nil
		}

		if err != nil {
			log.Printf("Error Receiving: %s", err)
			return err
		}

		log.Printf("Contact to store: %v", contact)

		server.ContactRepository.StoreContact(contact)

	}
}

func (server *ContactsServer) ListContacts(filterContact *contactpb.Contact, stream contactpb.Contacts_ListContactsServer) error {

	contact, err := server.ContactRepository.FindContact(filterContact.Id)

	if err != nil {
		return err
	}

	log.Printf("Contact found: %v", contact)

	err = stream.Send(contact)

	if err != nil {
		return err
	}

	return io.EOF
}
