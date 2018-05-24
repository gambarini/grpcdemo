package dbutils

import (
	"gopkg.in/mgo.v2"
	"fmt"
)

type DB interface {

}

func DialMongoDB() (session *mgo.Session, err error){

	session, err = mgo.Dial("mongodb-set-0.mongodb-service,mongodb-set-1.mongodb-service,mongodb-set-2.mongodb-service")

	if err != nil {
		return session, fmt.Errorf("fail to dial to mongodb cluster, %s", err)
	}

	return session, err
}


