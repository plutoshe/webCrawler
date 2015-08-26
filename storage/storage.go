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

type StoreFormat struct {
	Url     string
	Content string
}

func Link2Db(dburl string) (*mgo.Session, error) {
	session, err := mgo.Dial(dburl)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

func Link2DbByDefault() (*mgo.Session, error) {
	return Link2Db(MONGODB_URL)
}

func Link2Collection(session *mgo.Session, dbname, username, password, collectionname string, auth bool) *mgo.Collection {
	mongoDb := session.DB(dbname)
	if auth {
		mongoDb.Login(username, password)
	}
	return mongoDb.C(collectionname)

}

func Link2CollectionByDefault(session *mgo.Session) *mgo.Collection {
	return Link2Collection(session, MONGODB_DB, MONGODB_USER, MONGODB_PWD, MONGODB_COLLECTION, true)
}

func StoreInsert(c *mgo.Collection, in StoreFormat) error {
	err := c.Insert(&in)
	return err
}
