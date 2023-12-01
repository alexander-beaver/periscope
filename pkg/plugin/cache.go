package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
)

type CacheEngine struct {
	cache *bigcache.BigCache
}

func (c *CacheEngine) Cache(key string, value []byte) bool {
	err := c.cache.Set(key, value)
	if err != nil {
		fmt.Println("Error caching value", err.Error())
		return false
	}
	fmt.Println("Cached value")
	return true

}

func (c *CacheEngine) Contains(key string) bool {
	entry, err := c.cache.Get(key)
	if err == nil && entry != nil {
		return true
	}
	return false
}
func (c *CacheEngine) GetCachedValue(key string) []byte {
	// Open file and return value
	entry, err := c.cache.Get(key)
	if err != nil {
		fmt.Println("Error getting cached value", err.Error())
		return nil
	}
	fmt.Println("Got cached value for key", key)
	return entry
}
func NewCache(duration time.Duration) *CacheEngine {
	c := CacheEngine{}

	conf := bigcache.Config{
		Shards:             1024,
		LifeWindow:         duration,
		CleanWindow:        duration * 2,
		MaxEntriesInWindow: 1000 * 10 * 60,
		HardMaxCacheSize:   8192,
	}
	cache, err := bigcache.New(context.Background(), conf)

	if err != nil {
		fmt.Println("Error creating cache", err.Error())
		return nil
	}
	c.cache = cache
	return &c
}
