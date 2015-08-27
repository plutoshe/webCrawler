package main

import (
	"testing"
)

func TestPattern(t *testing.T) {
	var testPattern = []struct {
		url   string
		match bool
	}{
		{"http://www.dianping.com/beijing/ddd", true},
		{"http://www.dianping.com/beijing/ddd/aaa", false},
		{"http://www.dianping.com/beijing/", false},
		{"http://www.dianping.com/beijing", false},
		{"http://www.dianping.com/beijing/ddd/eee/eee/eee", false},
	}
	for _, pattern := range testPattern {
		if checkMatchPattern("", pattern.url) != pattern.match {
			t.Fatalf(
				"Get unexpected pattern match result when match url (%s). Wanted = %t, Get = %t\n",
				pattern.url,
				pattern.match,
				!pattern.match,
			)
		}
	}
}
