package repetition

import (
	"testing"

	"gopkg.in/redis.v3"
)

func TestRepetiton(t *testing.T) {

	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rep := &RepetitionJudgement{}
	rep.InitializeVisited(c, "tt")
	rep.VisitedNewNode("www.baidu.com")
	rep.VisitedNewNode("www.dazhongdainping.com")
	var testdata = []struct {
		url     string
		visited bool
	}{
		{"www.baidu.com", true},
		{"123", false},
		{"www.dazhongdainping.com", true},
		{"", false},
		{"www.google.com", false},
	}
	for _, i := range testdata {
		ex, err := rep.CheckIfVisited(i.url)
		if err != nil {
			t.Fatalf("Repetition test failed, error Msg = %v", err)
		}
		if i.visited != ex {
			t.Fatalf("Repetition check failed, node visited status mismatch. Wanted = %t. Get = %t\n", i.visited, !i.visited)
		}
	}
}
