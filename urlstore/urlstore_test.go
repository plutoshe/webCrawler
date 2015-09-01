package urlstore

import (
	"testing"

	"gopkg.in/redis.v3"
)

func TestURLInit(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	c.Del("t1")
	c.Del("t2")
	c.Del("t3")
	c.SAdd("t1", "a", "b")
	c.SAdd("t2", "c", "d", "a")
	c.SAdd("t3", "1", "2", "3")

	var testdata = []struct {
		wantedNum int
		unionSet  []string
	}{
		{4, []string{"t1", "t2"}},
		{5, []string{"t1", "t3"}},
		{0, nil},
		{2, []string{"t1"}},
		{7, []string{"t1", "t2", "t3"}},
	}

	for _, i := range testdata {
		store := &URLCrawlerStore{}
		store.InitialURLsStore(c, "t5", "com", i.unionSet...)
		cmp, err := c.SCard("t5").Result()
		if err != nil {
			t.Fatalf("Test initialization of URLs storage failed, error message = %v\n", err)
		}
		if cmp != (int64)(i.wantedNum) {
			t.Fatalf("Test initialization of URLs storage failed. The storage size mismatch. Wanted = %d. Get = %d\n", i.wantedNum, cmp)
		}
	}

	return
}

func TestUploadURL(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var testdata = []struct {
		wantedNum int
		addString [][]string
	}{
		{5, [][]string{[]string{"t1", "t2"}, []string{"t1", "t3"}, []string{"t4", "t5"}}},
		{0, nil},
		{1, [][]string{[]string{"1", "1"}, []string{"1", "1"}, []string{"1"}}},
		{3, [][]string{[]string{"t1", "t2", "t3"}}},
	}

	for _, i := range testdata {
		store := &URLCrawlerStore{}
		store.InitialURLsStore(c, "t5", "com")
		for _, addition := range i.addString {
			_, err := store.UploadURL(addition...)
			if err != nil {
				t.Fatalf("Test Upload URL failed, error message = %v\n", err)
			}
		}
		cmp, err := c.SCard("t5").Result()
		if err != nil {
			t.Fatalf("Test Upload URL failed, error message = %v\n", err)
		}
		if cmp != (int64)(i.wantedNum) {
			t.Fatalf("Test Upload URL failed. The storage size mismatch. Wanted = %d. Get = %d\n", i.wantedNum, cmp)
		}
	}
}

func TestGetOneNeedCrawlerURL(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var testdata = []struct {
		wantedNum int
		addString [][]string
	}{
		{5, [][]string{[]string{"t1", "t2"}, []string{"t1", "t3"}, []string{"t4", "t5"}}},
		{3, [][]string{[]string{"t1", "t2", "t3"}}},
	}
	for _, i := range testdata {
		store := &URLCrawlerStore{}
		store.InitialURLsStore(c, "t5", "com")
		exist := make(map[string]bool)
		for _, addition := range i.addString {
			for _, node := range addition {
				exist[node] = true
			}
			_, err := store.UploadURL(addition...)
			if err != nil {
				t.Fatalf("Test Upload URL failed, error message = %v\n", err)
			}
		}
		for {
			url, err := store.GetOneNeedCrawlerURL()
			if url == "" {
				break
			}
			if err != nil {
				t.Fatalf("Test Get One Need Crawler URL failed, error message = %v\n", err)
			}
			_, ok := exist[url]
			if !ok {
				t.Fatalf("Test Get One Need Crawler URL failed. Get unwanted URL = %s.\n", url)
			}
		}
	}
}
