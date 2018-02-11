package config

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
)

var session *mgo.Session
var Db string
var Col string

func InitDb() {

	var err error
	var urlConection string

	MONGO_HOST := os.Getenv("MONGO_HOST")
	MONGO_USER := os.Getenv("MONGO_USER")
	MONGO_PASS := os.Getenv("MONGO_PASS")
	MONGO_PORT := os.Getenv("MONGO_PORT")
	MONGO_DATABASE := os.Getenv("MONGO_DATABASE")

	if len(MONGO_USER) > 0 && len(MONGO_PASS) > 0 {
		urlConection = fmt.Sprintf("mongodb://%s:%s@%s:%s", MONGO_USER, MONGO_PASS, MONGO_HOST, MONGO_PORT)
	} else {
		urlConection = fmt.Sprintf("mongodb://%s:%s/%s", MONGO_HOST, MONGO_PORT, MONGO_DATABASE)
	}

	session, err = mgo.Dial(urlConection)

	if err != nil {
		fmt.Println(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return
}

func Sess() *mgo.Session {
	return session
}

func Collection() (col *mgo.Collection, err error) {
	col = session.DB(Db).C(Col)
	return col, err
}
