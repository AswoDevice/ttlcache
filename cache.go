package ttlcache

import (
	"sync"
	"time"
)

// Cache is a synchronised map of items that auto-expire once stale
type Cache struct {
	mutex        sync.RWMutex
	ttl          time.Duration
	hasTouchLife bool
	items        map[string]*Item
}

// Set is a thread-safe way to add new items to the map
func (cache *Cache) Set(key string, data interface{}) {
	cache.mutex.Lock()
	item := &Item{data: data}
	item.touch(cache.ttl)
	cache.items[key] = item
	cache.mutex.Unlock()
}

// Get is a thread-safe way to lookup items
func (cache *Cache) Get(key string) (data interface{}, found bool) {
	cache.mutex.Lock()
	item, exists := cache.items[key]
	if !exists || item.expired() {
		data = ""
		found = false
	} else {
		if cache.hasTouchLife {
			item.touch(cache.ttl)
		}
		data = item.data
		found = true
	}
	cache.mutex.Unlock()
	return
}

// Delete is a thread-safe way to delete items to the map
func (cache *Cache) Delete(key string) {
	cache.mutex.Lock()
	delete(cache.items, key)
	cache.mutex.Unlock()
}

// Count returns the number of items in the cache
// (helpful for tracking memory leaks)
func (cache *Cache) Count() int {
	cache.mutex.RLock()
	count := len(cache.items)
	cache.mutex.RUnlock()
	return count
}

func (cache *Cache) cleanup() {
	cache.mutex.Lock()
	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
	cache.mutex.Unlock()
}

func (cache *Cache) startCleanupTimer() {
	duration := cache.ttl
	if duration < time.Second {
		duration = time.Second
	}
	ticker := time.Tick(duration)
	go (func() {
		for {
			select {
			case <-ticker:
				cache.cleanup()
			}
		}
	})()
}

// NewCache is a helper to create instance of the Cache struct
func NewCache(config Config) *Cache {
	cache := &Cache{
		ttl:          config.Duration,
		hasTouchLife: config.HasTouchLife,
		items:        map[string]*Item{},
	}
	cache.startCleanupTimer()
	return cache
}

type Config struct {
	Duration time.Duration
	HasTouchLife bool // Every lookup, also touches the item, hence extending it's life
}
