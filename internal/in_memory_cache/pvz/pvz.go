package pvz

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

// cacheItem struct in_memory_cache cacheItem
type cacheItem struct {
	value      model.PVZ
	expiration time.Time
	created    time.Time
}

// Cache struct in_memory_cache
type Cache struct {
	mu                sync.RWMutex
	items             map[uuid.UUID]cacheItem
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

// New Initializing a new memory in_memory_cache
func New(defaultExpiration time.Duration, cleanupInterval time.Duration) *Cache {
	cache := Cache{
		mu:                sync.RWMutex{},
		items:             make(map[uuid.UUID]cacheItem),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	cache.startCleaningExpiredItems()
	return &cache
}

// startCleaningExpiredItems starts cleaning
func (c *Cache) startCleaningExpiredItems() {
	go func() {
		for range time.Tick(c.cleanupInterval) {
			c.mu.Lock()
			for k, v := range c.items {
				if v.isExpired() {
					delete(c.items, k)
				}
			}
			c.mu.Unlock()
		}
	}()
}

// isExpired checks if the in_memory_cache item has expired
func (c cacheItem) isExpired() bool {
	return c.expiration.Before(time.Now())
}

// Set setting a in_memory_cache by key
func (c *Cache) Set(key uuid.UUID, value model.PVZ, ttl time.Duration) {
	if ttl == 0 {
		ttl = c.defaultExpiration
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
		created:    time.Now(),
	}
}

// Get getting a in_memory_cache by key
func (c *Cache) Get(key uuid.UUID) (model.PVZ, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return model.PVZ{}, false
	}

	if item.isExpired() {
		delete(c.items, key)
		return model.PVZ{}, false
	}

	return item.value, true
}

// Delete in_memory_cache by key
func (c *Cache) Delete(key uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}
