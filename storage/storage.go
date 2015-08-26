package storage

import "gopkg.in/mgo.v2"

// mongodb configuration
const (
	MONGODB_URL        = "127.0.0.1:27017"
	MONGODB_DB         = "test"
	MONGODB_USER       = "testuser1"
	MONGODB_PWD        = "qwe123123"
	MONGODB_COLLECTION = "urlcollection"
)

type storeFormat struct {
	Url     string
	Content string
}

func link2Db(dburl string) (*mgo.Session, error) {
	session, err := mgo.Dial(dburl)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

func link2DbByDefault() (*mgo.Session, error) {
	return link2Db(MONGODB_URL)
}

func link2Collection(session *mgo.Session, dbname, username, password, collectionname string, auth bool) *mgo.Collection {
	mongoDb := session.DB(dbname)
	if auth {
		mongoDb.Login(username, password)
	}
	return mongoDb.C(collectionname)

}

func link2CollectionByDefault(session *mgo.Session) *mgo.Collection {
	return link2Collection(session, MONGODB_DB, MONGODB_USER, MONGODB_PWD, MONGODB_COLLECTION, true)
}

func storeInsert(c *mgo.Collection, in storeFormat) error {
	err := c.Insert(&in)
	return err
}
