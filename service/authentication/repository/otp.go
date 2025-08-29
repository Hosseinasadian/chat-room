package repository

import "github.com/hosseinasadian/chat-application/adapter/redis"

type OTP struct {
	adapter redis.Adapter
}

func New(adapter redis.Adapter) OTP {
	return OTP{adapter: adapter}
}

func (repo OTP) Adapter() redis.Adapter {
	return repo.adapter
}
