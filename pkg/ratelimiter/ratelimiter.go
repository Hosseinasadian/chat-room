package ratelimiter

import (
	"context"
	"fmt"
	"github.com/hosseinasadian/chat-application/pkg/richerror"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	limit  int64
	window time.Duration
}

func New(client *redis.Client, limit int64, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (l *RateLimiter) Allow(ctx context.Context, key string) error {
	const op = "rateLimiter.RateLimiter.Allow"

	redisKey := fmt.Sprintf("rate_limit:%s", key)

	pipe := l.client.Pipeline()
	val := pipe.Incr(ctx, redisKey)
	pipe.Expire(ctx, redisKey, l.window)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return richerror.New(op).WithWrapper(err).WithMessage("failed to check rate limit")
	}

	if val.Val() > l.limit {
		return richerror.New(op).WithKind(richerror.KindTooManyRequests).WithMessage(http.StatusText(http.StatusTooManyRequests))
	}

	return nil
}
