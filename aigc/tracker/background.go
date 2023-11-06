package tracker

import (
	"fmt"
	"sync"
)

const (
	defaultCacheCount = 64
	defaultCacheSize  = 16 * 1024 // 16KB
)

// Background send log or tracker on background
type Background struct {
	client *Client

	cacheArr []LogCache
	nowSize  int

	cacheCount int
	cacheSize  int

	lock *sync.Mutex
}

// NewBackground ...
func NewBackground(client *Client) *Background {
	return &Background{
		client: client,

		cacheArr: make([]LogCache, 0, defaultCacheCount),
		nowSize:  0,

		cacheCount: defaultCacheCount,
		cacheSize:  defaultCacheSize,

		lock: new(sync.Mutex),
	}
}

// Push a log or tracker
func (b *Background) Push(item LogCache) {
	go b.push(item)
}

func (b *Background) push(item LogCache) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.cacheArr = append(b.cacheArr, item)
	b.nowSize += item.Size()

	if len(b.cacheArr) >= b.cacheCount || b.nowSize >= b.cacheSize {
		b.sendAll()
	}
}

func (b *Background) sendAll() {
	logArr := make([]Log, 0, len(b.cacheArr))
	trackerArr := make([]Tracker, 0, len(b.cacheArr))

	for _, item := range b.cacheArr {
		switch itemWithType := item.(type) {
		case *Log:
			logArr = append(logArr, *itemWithType)
		case *Tracker:
			trackerArr = append(trackerArr, *itemWithType)
		default:
			fmt.Printf("%+v is not Log or Tracker", item)
		}
	}

	// 清理
	b.cacheArr = b.cacheArr[:0]
	b.nowSize = 0

	// 投递日志和埋点
	if len(logArr) != 0 {
		b.client.SendLogs(logArr)
	}
	if len(trackerArr) != 0 {
		b.client.SendTrackers(trackerArr)
	}
}

// SendAll log or tracker in cache
func (b *Background) SendAll() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.sendAll()
}
