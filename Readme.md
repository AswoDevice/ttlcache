## TTLCache - an in-memory LRU cache with expiration

TTLCache is a minimal wrapper over a string map in golang, entries of which are

1. Thread-safe
2. Auto-Expiring after a certain time
3. Auto-Extending expiration on `Get`s

[![Build Status](https://travis-ci.org/wunderlist/ttlcache.svg)](https://travis-ci.org/wunderlist/ttlcache)

#### Usage
```go
import (
  "time"
  "github.com/wunderlist/ttlcache"
)

func main () {
  cache := ttlcache.NewCache(ttlcache.Config{
  	Duration: time.Minute * 8,
  	HasTouchLife: true,
  })
  cache.Set("key", "value")
  value, exists := cache.GetString("key")
  count := cache.Count()
  cache.Delete("key")
}
```