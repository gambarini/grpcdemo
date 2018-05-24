package db

import (
	"gopkg.in/mgo.v2"
)


type DB struct {
	Session *mgo.Session
}

func NewDB(session *mgo.Session) *DB {
	return &DB{session}

}


