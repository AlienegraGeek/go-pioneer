package util

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	cacheInstance *cache.Cache
	cacheOnce     sync.Once
)

func GetCacheInstance() *cache.Cache {
	cacheOnce.Do(func() {
		// 创建一个缓存实例，设置默认过期时间为 5 分钟，清理过期缓存的周期为 10 分钟
		cacheInstance = cache.New(10*time.Minute, 20*time.Minute)
	})
	return cacheInstance
}
