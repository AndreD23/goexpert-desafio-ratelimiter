package limiter

import (
	"context"
	"github.com/AndreD23/goexpert-desafio-ratelimiter/internal/redis"
	"time"
)

type RateLimiter struct {
	redisStore    redis.RedisInterface
	requestsIP    int
	requestsToken int
	blockDuration time.Duration
}

func NewRateLimiter(r redis.RedisInterface, reqIP, reqToken int, blockDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		redisStore:    r,
		requestsIP:    reqIP,
		requestsToken: reqToken,
		blockDuration: blockDuration,
	}
}

func (rl *RateLimiter) CheckRateLimitIP(ip string) (bool, error) {
	isBlocked, err := rl.redisStore.IsBlocked(context.Background(), ip)
	if err != nil {
		return false, err
	}
	if isBlocked {
		return true, nil
	}

	requestTimes, err := rl.redisStore.GetRequestTimestamps(context.Background(), ip)
	if err != nil {
		return false, err
	}

	currentTimestamp := time.Now().UnixNano()
	oneSecondAgo := currentTimestamp - int64(time.Second)
	validRequests := filterValidRequests(requestTimes, oneSecondAgo)

	if len(validRequests) >= rl.requestsIP {
		err := rl.redisStore.BlockKey(context.Background(), ip, rl.blockDuration)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	err = rl.redisStore.AddRequestTimestamp(context.Background(), ip, currentTimestamp)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (rl *RateLimiter) CheckRateLimitToken(token string) (bool, error) {
	isBlocked, err := rl.redisStore.IsBlocked(context.Background(), token)
	if err != nil {
		return false, err
	}
	if isBlocked {
		return true, nil
	}

	requestTimes, err := rl.redisStore.GetRequestTimestamps(context.Background(), token)
	if err != nil {
		return false, err
	}

	currentTimestamp := time.Now().UnixNano()
	oneSecondAgo := currentTimestamp - int64(time.Second)
	validRequests := filterValidRequests(requestTimes, oneSecondAgo)

	if len(validRequests) >= rl.requestsToken {
		err := rl.redisStore.BlockKey(context.Background(), token, rl.blockDuration)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	err = rl.redisStore.AddRequestTimestamp(context.Background(), token, currentTimestamp)
	if err != nil {
		return false, err
	}

	return false, nil
}

func filterValidRequests(requestTimes []int64, threshold int64) []int64 {
	var validRequests []int64
	for _, t := range requestTimes {
		if t > threshold {
			validRequests = append(validRequests, t)
		}
	}
	return validRequests
}
