package cache

import (
	"github.com/astaxie/beego/cache"
)

func NewCache() cache.Cache {
	return cache.NewMemoryCache()
}