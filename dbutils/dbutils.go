package dbutils

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"time"
	"log"
)

func DialMongoDB() (session *mgo.Session, err error){

	session, err = mgo.Dial("mongodb-set-0.mongodb-service,mongodb-set-1.mongodb-service,mongodb-set-2.mongodb-service")

	if err != nil {
		return session, fmt.Errorf("fail to dial to mongodb cluster, %s", err)
	}

	return session, err
}

type DB struct {
	Session *mgo.Session
	Ticker *time.Ticker
}

func NewDB(session *mgo.Session) *DB {

	ticker := time.NewTicker(time.Minute * 5)

	db := &DB{session, ticker}

	go db.SessionRefresh()

	return db
}

func (db *DB) SessionRefresh() {

	for {
		<- db.Ticker.C

		log.Print("Refreshing mongodb session...")
		db.Session.Refresh()
	}
}


