package dbutils

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"time"
	"log"
)

const (
	MongoDBURL = "mongodb-set-0.mongodb-service,mongodb-set-1.mongodb-service,mongodb-set-2.mongodb-service"
)

type (

	DB interface {
		GetSession() *mgo.Session
		CleanUp()
	}

	MongoDB struct {
		session *mgo.Session
		ticker *time.Ticker
	}
)


func NewMongoDB(mongoDbURL string) (db *MongoDB, err error) {

	session, err := mgo.Dial(mongoDbURL)

	if err != nil {
		return db, fmt.Errorf("fail to dial to mongodb cluster, %s", err)
	}

	ticker := time.NewTicker(time.Minute * 10)

	db = &MongoDB{session, ticker}

	go db.SessionRefresh()

	return db, nil
}

func (db *MongoDB) SessionRefresh() {

	for {
		<- db.ticker.C

		log.Print("Refreshing mongodb session...")
		db.session.Refresh()
	}
}

func (db *MongoDB) GetSession() *mgo.Session{

	session := db.session.Copy()

	session.SetMode(mgo.Monotonic, true)

	return session
}

func (db *MongoDB) CleanUp() {
	db.ticker.Stop()
	db.session.Close()
}



