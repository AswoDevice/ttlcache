package ttlcache

func (cache *Cache) GetString(key string) (string, bool) {
	if data, found := cache.Get(key); found {
		return data.(string), true
	}

	return "", false
}

func (cache *Cache) GetInt(key string) (int, bool) {
	if data, found := cache.Get(key); found {
		return data.(int), true
	}

	return 0, false
}

func (cache *Cache) GetBytes(key string) ([]byte, bool) {
	if data, found := cache.Get(key); found {
		return data.([]byte), true
	}

	return nil, false
}
