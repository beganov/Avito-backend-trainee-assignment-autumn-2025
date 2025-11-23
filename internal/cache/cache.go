package cache

import (
	"sync"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
)

var UserCache lruCache
var TeamCache lruCache
var PRcache lruCache

// node in LRU list
type lruNode struct {
	key   string
	value interface{}
	prev  *lruNode
	next  *lruNode
}

// simple LRU cache for orders
type lruCache struct {
	capacity int
	store    map[string]*lruNode
	head     *lruNode
	tail     *lruNode
	mu       sync.Mutex
}

// constructor
func NewOrderCache(cap int) lruCache {
	return lruCache{
		capacity: cap,
		store:    make(map[string]*lruNode),
	}
}

// move node to front (most recently used)
func (c *lruCache) moveToFront(node *lruNode) {
	if c.head == node {
		return
	}

	// unlink node
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if c.tail == node {
		c.tail = node.prev
	}

	// put node at head
	node.prev = nil
	node.next = c.head
	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

	if c.tail == nil {
		c.tail = node
	}
}

// add new val to cache
func (c *lruCache) Set(key string, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node := &lruNode{key: key, value: val}
	c.store[key] = node
	c.moveToFront(node)

	if len(c.store) > c.capacity {
		delete(c.store, c.tail.key)
		if c.tail.prev != nil {
			c.tail = c.tail.prev
			c.tail.next = nil
		} else {
			c.head = nil
			c.tail = nil
		}
	}
}

// get val from cache
func (c *lruCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.store[key]; ok {
		c.moveToFront(node)
		return node.value, true
	}
	return nil, false
}

func InitCache() {
	UserCache = NewOrderCache(config.CacheCap)
	TeamCache = NewOrderCache(config.CacheCap)
	PRcache = NewOrderCache(config.CacheCap)
}
