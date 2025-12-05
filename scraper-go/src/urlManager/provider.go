package urlManager

import (
    "container/heap"
    "scraper/src/customTypes"
    "sync"
)

// PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*customTypes.StoreUrl

// UrlQueue encapsulates the priority queue state to allow multiple instances.
type UrlQueue struct {
    links          PriorityQueue
    mutex          sync.Mutex
    seen           map[string]struct{}
    processCounter int
    processMutex   sync.Mutex
}

// NewUrlQueue creates and initializes a new isolated queue instance.
func NewUrlQueue() *UrlQueue {
    uq := &UrlQueue{
        links: make(PriorityQueue, 0),
        seen:  make(map[string]struct{}),
    }
    heap.Init(&uq.links)
    return uq
}

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
func (pq *PriorityQueue) update(item *customTypes.StoreUrl, url string, priority int) {
    item.Url = url
    item.Priority = priority
    heap.Fix(pq, item.Index)
}

// AddUrl adds a URL to the specific queue instance.
func (uq *UrlQueue) AddUrl(url string, query string, level int) bool {
    uq.mutex.Lock()
    defer uq.mutex.Unlock()

    if _, ok := uq.seen[url]; ok {
        return false
    }

    isValid, validUrl := IsUrlValid(url, "")
    if !isValid {
        return false
    }

    uq.seen[url] = struct{}{}

    urlItem := &customTypes.StoreUrl{Priority: RankUrl(url, query, level), Url: validUrl}
    heap.Push(&uq.links, urlItem)

    return true
}

// GetUrl retrieves the highest priority URL from the specific queue instance.
func (uq *UrlQueue) GetUrl() *customTypes.StoreUrl {
    uq.mutex.Lock()
    defer uq.mutex.Unlock()

    if uq.links.Len() == 0 {
        return nil
    }

    link := heap.Pop(&uq.links).(*customTypes.StoreUrl)
    return link
}

// GetQueueLength returns the current size of the queue.
func (uq *UrlQueue) GetQueueLength() int {
    uq.mutex.Lock()
    defer uq.mutex.Unlock()
    return uq.links.Len()
}

// IncrementProcessCounter increments the processed count for this session.
func (uq *UrlQueue) IncrementProcessCounter() {
    uq.processMutex.Lock()
    defer uq.processMutex.Unlock()
    uq.processCounter++
}

// GetProcessedUrlsCount returns the processed count for this session.
func (uq *UrlQueue) GetProcessedUrlsCount() int {
    uq.processMutex.Lock()
    defer uq.processMutex.Unlock()
    return uq.processCounter
}

// ClearQueue resets the specific queue instance.
func (uq *UrlQueue) ClearQueue() {
    uq.mutex.Lock()
    defer uq.mutex.Unlock()
    uq.links = make(PriorityQueue, 0)
    heap.Init(&uq.links)
	uq.seen = make(map[string]struct{})
}