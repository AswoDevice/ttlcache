package ttlcache

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	cache := NewCache(Config{Duration: time.Minute})

	data, exists := cache.Get("hello")
	if exists || data != "" {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("hello", "world")
	data, exists = cache.Get("hello")
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if data != "world" {
		t.Errorf("Expected cache to return `world` for `hello`")
	}
}

func TestDelete(t *testing.T) {
	cache := NewCache(Config{Duration: time.Minute})

	cache.Set("hello", "world")
	_, exists := cache.GetString("hello")
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}

	cache.Delete("hello")

	_, exists = cache.GetString("hello")
	if exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
}

func TestExpiration(t *testing.T) {
	cache := NewCache(Config{Duration: time.Second, HasTouchLife: true})

	cache.Set("x", "1")
	cache.Set("y", "z")
	cache.Set("z", "3")

	count := cache.Count()
	if count != 3 {
		t.Errorf("Expected cache to contain 3 items")
	}

	<-time.After(500 * time.Millisecond)
	cache.mutex.Lock()
	cache.items["y"].touch(time.Second)
	item, exists := cache.items["x"]
	cache.mutex.Unlock()
	if !exists || item.data != "1" || item.expired() {
		t.Errorf("Expected `x` to not have expired after 200ms")
	}

	<-time.After(time.Second)
	cache.mutex.RLock()
	_, exists = cache.items["x"]
	if exists {
		t.Errorf("Expected `x` to have expired")
	}
	_, exists = cache.items["z"]
	if exists {
		t.Errorf("Expected `z` to have expired")
	}
	_, exists = cache.items["y"]
	if !exists {
		t.Errorf("Expected `y` to not have expired")
	}
	cache.mutex.RUnlock()

	count = cache.Count()
	if count != 1 {
		t.Errorf("Expected cache to contain 1 item")
	}

	<-time.After(600 * time.Millisecond)
	cache.mutex.RLock()
	_, exists = cache.items["y"]
	if exists {
		t.Errorf("Expected `y` to have expired")
	}
	cache.mutex.RUnlock()

	count = cache.Count()
	if count != 0 {
		t.Errorf("Expected cache to be empty")
	}
}

func TestHasTouchLife(t *testing.T) {
	cache := NewCache(Config{Duration: time.Second, HasTouchLife: false})

	cache.Set("hello", "world")

	<-time.After(500 * time.Millisecond)

	_, exists := cache.Get("hello")
	if !exists {
		t.Errorf("Expected `hello` to have expired after 500ms")
	}

	<-time.After(600 * time.Millisecond)

	_, exists = cache.Get("hello")
	if exists {
		t.Errorf("Expected `hello` to not have expired after 1100ms")
	}
}