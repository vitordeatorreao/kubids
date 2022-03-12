package kubid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const kubidTTL = 5 * time.Second

const maxCount = 1<<12 - 1

type CollisionCounter interface {
	GetCollisionCount(key string) (int64, error)
}

type redisCollisionCounter struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisRandClient(ctx context.Context, addr string, pw string, db int) CollisionCounter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
	return &redisCollisionCounter{rdb: rdb, ctx: ctx}
}

func (rr *redisCollisionCounter) GetCollisionCount(key string) (int64, error) {
	fkey := fmt.Sprintf("kubid:%s", key)
	pipe := rr.rdb.Pipeline()
	set := pipe.SetNX(rr.ctx, fkey, -1, kubidTTL)
	incr := pipe.Incr(rr.ctx, fkey)
	pipe.Expire(rr.ctx, fkey, kubidTTL)
	_, pipeErr := pipe.Exec(rr.ctx)
	if pipeErr != nil {
		return -1, pipeErr
	}
	if err := set.Err(); err != nil {
		return -1, err
	}
	if err := incr.Err(); err != nil {
		return -1, err
	}
	count := incr.Val()
	if count > maxCount {
		return -1, errors.New("counter overflow, try again next millisecond")
	}
	return count, nil
}
