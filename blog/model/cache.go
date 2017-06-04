package model

import (
	"fmt"
	"log"
	"time"

	"github.com/toolkits/cache"
	"github.com/ulricqin/blog/config"
)

func InitCache() {
	c := config.G
	if c.Cache.Provider == "memory" {
		cache.Instance = cache.NewInMemoryCache(c.Cache.Expire)
	} else if c.Cache.Provider == "redis" {
		cache.Instance = cache.NewRedisCache(
			c.Redis.Addr,
			c.Redis.Idle,
			c.Redis.Max,
			time.Duration(c.Redis.Timeout.Conn)*time.Millisecond,
			time.Duration(c.Redis.Timeout.Read)*time.Millisecond,
			time.Duration(c.Redis.Timeout.Write)*time.Millisecond,
			c.Cache.Expire,
		)
	} else {
		log.Fatalln("cache.provider should be memory or redis")
	}
}

func CacheKeyArticle(id int64) string {
	return fmt.Sprintf("/cache/article/%d", id)
}

func CacheKeyArticleId(url string) string {
	return fmt.Sprintf("/cache/articleid/%s", url)
}

func CacheKeyArticleTotal() string {
	return "/cache/article/total"
}

func CacheKeyPage(name string) string {
	return fmt.Sprintf("/cache/page/%s", name)
}
