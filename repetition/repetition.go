// This package check whether the url is visited or not.
package repetition

import "gopkg.in/redis.v3"

type RepetitionJudgement struct {
	redisServer *redis.Client
	collection  string
}

func (t *RepetitionJudgement) InitializeVisited(c *redis.Client, col string) error {
	t.redisServer = c
	t.collection = col
	return t.redisServer.Del(t.collection).Err()
}

func (t *RepetitionJudgement) VisitedNewNode(value ...string) (int64, error) {
	return t.redisServer.SAdd(t.collection, value...).Result()
}

func (t *RepetitionJudgement) CheckIfVisited(value string) (bool, error) {
	return t.redisServer.SIsMember(t.collection, value).Result()
}
