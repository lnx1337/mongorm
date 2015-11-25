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

	MONGO_HOST := os.Getenv("MONGO_HOST")
	MONGO_USER := os.Getenv("MONGO_USER")
	MONGO_PASS := os.Getenv("MONGO_PASS")
	MONGO_PORT := os.Getenv("MONGO_PORT")

	urlConection := fmt.Sprintf("mongodb://%s:%s@%s:%s", MONGO_USER, MONGO_PASS, MONGO_HOST, MONGO_PORT)

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