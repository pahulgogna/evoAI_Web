package urlManager

import (
	"container/heap"
	"scraper/src/customTypes"
	"sync"
)

type PriorityQueue []*customTypes.StoreUrl

var (
	Links       PriorityQueue
	mutex       sync.Mutex
	initialized bool
)

var (
	seen = make(map[string]struct{})
)

var (
	processCounter int
	processMutex   sync.Mutex
)

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Higher priority items come first
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*customTypes.StoreUrl)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // allow GC to reclaim the item
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
// Must be called while holding the mutex.
func (pq *PriorityQueue) update(item *customTypes.StoreUrl, url string, priority int) {
	item.Url = url
	item.Priority = priority
	heap.Fix(pq, item.Index)
}

// initHeap initializes the heap if not already initialized.
// Must be called while holding the mutex.
func initHeap() {
	if !initialized {
		Links = make(PriorityQueue, 0)
		heap.Init(&Links)
		initialized = true
	}
}

func AddUrl(url string, query string, level int) bool {
	mutex.Lock()
	defer mutex.Unlock()

	initHeap()
	
	if _, ok := seen[url]; ok { 
		return false
	}
	
	isValid, validUrl := IsUrlValid(url, "")
	if !isValid {
		return false
	}

	seen[url] = struct{}{}

	urlItem := &customTypes.StoreUrl{Priority: RankUrl(url, query, level), Url: validUrl}
	heap.Push(&Links, urlItem)

	return true
}

func GetUrl() *customTypes.StoreUrl {

	mutex.Lock()
	defer mutex.Unlock()

	if !initialized || Links.Len() == 0 {
		return nil
	}

	link := heap.Pop(&Links).(*customTypes.StoreUrl)
	return link
}

func GetQueueLength() int {
	mutex.Lock()
	defer mutex.Unlock()
	return Links.Len()
}

func IncrementProcessCounter() {
	processMutex.Lock()
	defer processMutex.Unlock()
	processCounter++
}

func GetProcessedUrlsCount() int {
	processMutex.Lock()
	defer processMutex.Unlock()
	return processCounter
}

func ClearQueue() {
	mutex.Lock()
	defer mutex.Unlock()
	Links = make(PriorityQueue, 0)
	heap.Init(&Links)
	seen = make(map[string]struct{})
	initialized = true
}