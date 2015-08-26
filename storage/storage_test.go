package storage

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestStore2Mongo(t *testing.T) {
	dbSession, err := link2DbByDefault()
	defer dbSession.Close()
	if err != nil {
		t.Fatal(err)
	}
	c := link2CollectionByDefault(dbSession)
	insertSet := storeFormat{"http://www.baidu.com", "only testing"}
	fmt.Println("store in mongodb")
	fmt.Println("url : ", insertSet.Url)
	fmt.Println("content : ", insertSet.Content)
	fmt.Println("=============")
	err = storeInsert(c, insertSet)
	if err != nil {
		t.Fatal(err)
	}
	result := storeFormat{}
	err = c.Find(bson.M{"url": "http://www.baidu.com"}).One(&result)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)
}
