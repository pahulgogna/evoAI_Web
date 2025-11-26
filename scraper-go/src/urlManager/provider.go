package urlManager

import (
	// "fmt"
	"sync"
)


var (
	Links []string = []string{}
	linksMutex sync.Mutex
)

var (
	Seen map[string]int8 = map[string]int8{}
	MapMutex sync.Mutex
)

var (
	processCounter int
	processeMutex sync.Mutex
)


func AddUrl(url string) bool {

	MapMutex.Lock()
	if _, ok := Seen[url]; ok {
		MapMutex.Unlock()
		return !ok
	}

	Seen[url] = 1
	MapMutex.Unlock()

	linksMutex.Lock()
	defer linksMutex.Unlock()

	isValid, url := IsUrlValid(url, "")
	if !isValid {
		return false
	}

	Links = append(Links, url)

	return true
}

func GetUrl() string {
	linksMutex.Lock()
	defer linksMutex.Unlock()
	
	if len(Links) == 0 {
		return ""
	}
	
	link := Links[len(Links) - 1]
	
	Links = Links[:len(Links) - 1]

	return link
}

func IncrementProcessCounter() {
	processeMutex.Lock()
	defer processeMutex.Unlock()

	processCounter = processCounter + 1
}

func GetProcessedUrlsCount() int {
	return processCounter
}