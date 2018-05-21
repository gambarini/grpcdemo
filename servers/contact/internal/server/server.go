package server

import pb "github.com/gambarini/grpcdemo/pb/contact"
import (
	"github.com/gambarini/grpcdemo/servers/contact/internal/db"
	"io"
	"log"
)

type ContactsServer struct {
	DB *db.DB
}

func (server *ContactsServer) StoreContacts(stream pb.Contacts_StoreContactsServer) error {

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

		server.DB.StoreContact(contact)

	}
}

func (server *ContactsServer) ListContacts(filterContact *pb.Contact, stream pb.Contacts_ListContactsServer) error {

	contact, err := server.DB.FindContact(filterContact.Id)

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
