package cacheModel

import (
	"github.com/zeromicro/go-zero/core/collection"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/key/cache"
	"time"
)

// 内存缓存

type IMMemoryCache struct {
	c       *collection.Cache
	userRpc usercenter.Usercenter
}

func NewIMMemoryCache(rpc usercenter.Usercenter) *IMMemoryCache {
	c, err := collection.NewCache(cache.DefaultMemoryCacheExpireMinute*time.Minute, collection.WithLimit(cache.MaxMemoryCacheCnt))
	if err != nil {
		panic(err)
	}
	return &IMMemoryCache{c: c, userRpc: rpc}
}
