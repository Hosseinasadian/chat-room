package repository

import "github.com/hosseinasadian/chat-application/adapter/redis"

type Cache struct {
	adapter redis.Adapter
}

func New(adapter redis.Adapter) Cache {
	return Cache{adapter: adapter}
}

func (Cache Cache) Adapter() redis.Adapter {
	return Cache.adapter
}
