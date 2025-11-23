package cache

import (
	"sync"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
)

var UserCache lruCache
var TeamCache lruCache
var PRcache lruCache

type lruNode struct {
	key   string
	value interface{}
	prev  *lruNode
	next  *lruNode
}

type lruCache struct {
	capacity int
	store    map[string]*lruNode
	head     *lruNode
	tail     *lruNode
	mu       sync.Mutex
}

func NewOrderCache(cap int) lruCache {
	return lruCache{
		capacity: cap,
		store:    make(map[string]*lruNode),
	}
}

func (c *lruCache) moveToFront(node *lruNode) {
	if c.head == node {
		return
	}

	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if c.tail == node {
		c.tail = node.prev
	}

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
