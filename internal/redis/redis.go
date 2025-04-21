package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) IncrementRequestCount(ctx context.Context, key string, duration time.Duration) (int, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	_, err = r.client.Expire(ctx, key, duration).Result()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *Redis) GetRequestCount(ctx context.Context, key string) (int, error) {
	count, err := r.client.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Redis) BlockKey(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Set(ctx, key+":blocked", true, duration).Err()
}

func (r *Redis) IsBlocked(ctx context.Context, key string) (bool, error) {
	blocked, err := r.client.Get(ctx, key+":blocked").Bool()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return blocked, nil
}

func (r *Redis) SetRequestTimestamp(ctx context.Context, key string) error {
	timestamp := time.Now().UnixNano()
	return r.client.Set(ctx, key+":timestamp", timestamp, 0).Err()
}

func (r *Redis) GetRequestTimestamp(ctx context.Context, key string) (int64, error) {
	timestamp, err := r.client.Get(ctx, key+":timestamp").Int64()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return timestamp, nil
}

func (r *Redis) AddRequestTimestamp(ctx context.Context, key string, timestamp int64) error {
	return r.client.ZAdd(ctx, key+":timestamps", &redis.Z{
		Score:  float64(timestamp),
		Member: timestamp,
	}).Err()
}

func (r *Redis) GetRequestTimestamps(ctx context.Context, key string) ([]int64, error) {
	timestamps, err := r.client.ZRangeWithScores(ctx, key+":timestamps", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var result []int64
	for _, ts := range timestamps {
		result = append(result, int64(ts.Score))
	}

	return result, nil
}

func (r *Redis) CleanupOldTimestamps(ctx context.Context, key string, threshold int64) error {
	return r.client.ZRemRangeByScore(ctx, key+":timestamps", "0", fmt.Sprintf("%d", threshold)).Err()
}
