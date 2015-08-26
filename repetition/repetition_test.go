package repetition

import (
	"testing"
)

func TestRepetiton(t *testing.T) {
	InitializeVisited()
	VisitedNewNode("www.baidu.com")
	if !checkIfVisited("www.baidu.com") {
		t.Fatal("Repetition check failed : supposed visited node shows non-visited.")
	}
	if CheckIfVisited("www.google.com") {
		t.Fatal("Repetition check failed : supposed non-visited node shows visited.")
	}
}
