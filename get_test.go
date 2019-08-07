package ttlcache

import (
	"testing"
	"time"
)

func TestGetString(t *testing.T) {
	cache := &Cache{
		ttl:   time.Minute,
		items: map[string]*Item{},
	}

	data, exists := cache.GetString("hello")
	if exists || data != "" {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("hello", "world")
	data, exists = cache.GetString("hello")
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if data != "world" {
		t.Errorf("Expected cache to return `world` for `hello`")
	}
}

func TestGetBytes(t *testing.T) {
	cache := &Cache{
		ttl:   time.Minute,
		items: map[string]*Item{},
	}

	data, exists := cache.GetBytes("hello")
	if exists || data != nil {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("hello", []byte("world"))
	data, exists = cache.GetBytes("hello")
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if string(data) != "world" {
		t.Errorf("Expected cache to return `world` for `hello`")
	}
}

func TestGetInt(t *testing.T) {
	cache := &Cache{
		ttl:   time.Minute,
		items: map[string]*Item{},
	}

	data, exists := cache.GetInt("hello")
	if exists || data != 0 {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("hello", 123)
	data, exists = cache.GetInt("hello")
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if data != 123 {
		t.Errorf("Expected cache to return `123` for `hello`")
	}
}
