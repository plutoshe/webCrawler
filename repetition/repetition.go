// This package check whether the url is visited or not.
package repetition

var visited map[string]bool

func initializeVisited() {
	visited = make(map[string]bool)
}

func visitedNewNode(key string) {
	visited[key] = true
}

func checkIfVisited(key string) bool {
	_, ok := visited[key]
	return ok
}
