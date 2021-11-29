package cache

import "time"
import gocache "github.com/patrickmn/go-cache"

const DefaultMaxLimit = 10000

var (
	gMemoryCache *gocache.Cache
	gEmptyCache  *gocache.Cache
)

func InitMemoryCache(expire time.Duration, clean time.Duration) {
	gMemoryCache = gocache.New(expire, clean)
	gEmptyCache = gocache.New(expire, clean)
}

func GetMemoryCache() *gocache.Cache {
	return gMemoryCache
}

func GetEmptyCache() *gocache.Cache {
	return gEmptyCache
}
