package cache

import "time"
import gocache "github.com/patrickmn/go-cache"

const DefaultMaxLimit = 10000

var (
	memoryCache *gocache.Cache
	emptyCache  *gocache.Cache
)

func InitMemoryCache(expire time.Duration, clean time.Duration) {
	memoryCache = gocache.New(expire, clean)
	emptyCache = gocache.New(expire, clean)
}

func GetMemoryCache() *gocache.Cache {
	return memoryCache
}

func GetEmptyCache() *gocache.Cache {
	return emptyCache
}
