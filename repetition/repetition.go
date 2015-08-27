// This package check whether the url is visited or not.
package repetition

var visited map[string]bool

func InitializeVisited() {
	visited = make(map[string]bool)
}

func VisitedNewNode(key string) {
	visited[key] = true
}

func CheckIfVisited(key string) bool {
	_, ok := visited[key]
	return ok
}
