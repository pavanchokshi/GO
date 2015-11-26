package components

import (
	"gopkg.in/mgo.v2"
)


func ConnectMongo() *mgo.Session {
	uri := "mongodb://dbuser:dbuser@ds041154.mongolab.com:41154/test_db"
    sess, err := mgo.Dial(uri)
//	defer sess.Close()
  	if err != nil {
    	panic(err)
  	} else {
	  	sess.SetSafe(&mgo.Safe{})
	    //collection = sess.DB("test_db").C("test")
	}
	return sess
}
