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
	cr := rr.rdb.SetNX(rr.ctx, fkey, rc, kubidTTL)
	if err := cr.Err(); err != nil {
		return 0, err
	}
	if cr.Val() {
		return rc, nil
	}
	ir := rr.rdb.Incr(rr.ctx, fkey)
	if err := ir.Err(); err != nil {
		return 0, err
	}
	rd := ir.Val()
	if rd > math.MaxUint32 {
		return 0, errors.New("counter overflow, try again next second")
	}
	return uint32(ir.Val()), nil
}
