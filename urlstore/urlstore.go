// This package implements a container storing the urls which should be
// crawled in db(redis currently)
//
// Therefore, the crawler could work correct concurrently in different
// thread/process/machine.
//
// The contianer stores the url need to crawl, removes the url alreay has
// been crawled.

package urlstore

import (
	"strings"

	"gopkg.in/redis.v3"
)

type URLCrawlerStore struct {
	redisServer          *redis.Client
	collectionNeedCrawl  string
	collectionNeedCommit string
}

func (t *URLCrawlerStore) InitialURLsStore(c *redis.Client, colNeedCrawl string, colNeedCommit string, cmb ...string) (int64, error) {
	t.redisServer = c
	t.collectionNeedCrawl = colNeedCrawl
	t.collectionNeedCommit = colNeedCommit

	if cmb == nil {
		t.redisServer.Del(t.collectionNeedCommit)
		return t.redisServer.Del(colNeedCrawl).Result()
	}

	result, err := t.redisServer.SUnionStore(t.collectionNeedCrawl, cmb...).Result()
	t.redisServer.Del(t.collectionNeedCommit)
	return result, err
}

func (t *URLCrawlerStore) GetOneNeedCrawlerURL() (string, error) {
	url, err := t.redisServer.SPop(t.collectionNeedCrawl).Result()
	if err != nil {
		if strings.Contains(err.Error(), "redis: nil") {
			return "", nil
		}
		return "", err
	}

	rep, err := t.redisServer.SAdd(t.collectionNeedCommit, url).Result()
	if rep == 0 {
		return "", nil
	}
	if err != nil {
		t.redisServer.SAdd(t.collectionNeedCrawl, url)
		return "", err
	}
	return url, nil
}

func (t *URLCrawlerStore) UploadURL(value ...string) (int64, error) {
	return t.redisServer.SAdd(t.collectionNeedCrawl, value...).Result()
}

func (t *URLCrawlerStore) CommitURL(value ...string) (int64, error) {
	return t.redisServer.SRem(t.collectionNeedCommit, value...).Result()
}
