package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type lockRepository struct {
	client *redis.Client
}

func NewLockRepository(client *redis.Client) *lockRepository {
	return &lockRepository{client: client}
}

func (r *lockRepository) Acquire(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	ok, err := r.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *lockRepository) Release(ctx context.Context, key string, value string) (bool, error) {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		end
		return 0
	`
	result, err := r.client.Eval(ctx, script, []string{key}, value).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

func (r *lockRepository) Extend(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("EXPIRE", KEYS[1], ARGV[2])
		end
		return 0
	`
	result, err := r.client.Eval(ctx, script, []string{key}, value, int(ttl.Seconds())).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}
