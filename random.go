package kubid

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
)

const kubidTTL = 5 * time.Second

type RandClient interface {
	SetOrGetRand(key string, rc uint32) (uint32, error)
}

type redisRandClient struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisRandClient(ctx context.Context, addr string, pw string, db int) RandClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
	return &redisRandClient{rdb: rdb, ctx: ctx}
}

func (rr *redisRandClient) SetOrGetRand(key string, rc uint32) (uint32, error) {
	fkey := fmt.Sprintf("kubid:%s", key)
	pipe := rr.rdb.Pipeline()
	set := pipe.SetNX(rr.ctx, fkey, rc, kubidTTL)
	incr := pipe.Incr(rr.ctx, fkey)
	pipe.Expire(rr.ctx, fkey, kubidTTL)
	_, pipeErr := pipe.Exec(rr.ctx)
	if pipeErr != nil {
		return 0, pipeErr
	}
	if err := set.Err(); err != nil {
		return 0, err
	}
	if err := incr.Err(); err != nil {
		return 0, err
	}
	rd := incr.Val()
	if rd > math.MaxUint32 {
		return 0, errors.New("counter overflow, try again next second")
	}
	return uint32(rd), nil
}
