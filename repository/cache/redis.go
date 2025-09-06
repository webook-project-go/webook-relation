package cache

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/webook-project-go/webook-relation/domain"
	"time"
)

type Cache interface {
	FindFollowers(ctx context.Context, uid int64) ([]int64, error)
	SetFollowers(ctx context.Context, uid int64, res []int64) error
	FindFollowees(ctx context.Context, uid int64) ([]int64, error)
	SetFollowees(ctx context.Context, uid int64, res []int64) error
	GetFollowerCount(ctx context.Context, uid int64) (uint32, error)
	SetFollowersCount(ctx context.Context, uid int64, cnt uint32) error
	UpdateRelation(ctx context.Context, info domain.RelationInfo, delta int) error
	GetFolloweeCount(ctx context.Context, uid int64) (uint32, error)
	SetFolloweesCount(ctx context.Context, uid int64, cnt uint32) error
	CacheLimit() int
}
type redisCache struct {
	rdb      redis.Cmdable
	listTTL  time.Duration
	countTTL time.Duration
}

//go:embed updateCount.lua
var updateRelationLua string

func New(rdb redis.Cmdable) Cache {
	return &redisCache{
		rdb:      rdb,
		listTTL:  time.Minute * 10,
		countTTL: time.Minute * 10,
	}
}

func followersKey(uid int64) string   { return fmt.Sprintf("followers:%d", uid) }
func followeesKey(uid int64) string   { return fmt.Sprintf("followees:%d", uid) }
func followCountKey(uid int64) string { return fmt.Sprintf("follow_count:%d", uid) }
func (c *redisCache) CacheLimit() int {
	return 1000
}
func (c *redisCache) FindFollowers(ctx context.Context, uid int64) ([]int64, error) {
	val, err := c.rdb.Get(ctx, followersKey(uid)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, redis.Nil
	} else if err != nil {
		return nil, err
	}
	var res []int64
	err = json.Unmarshal([]byte(val), &res)
	return res, err
}

func (c *redisCache) SetFollowers(ctx context.Context, uid int64, res []int64) error {
	data, _ := json.Marshal(res)
	return c.rdb.Set(ctx, followersKey(uid), data, c.listTTL).Err()
}

func (c *redisCache) FindFollowees(ctx context.Context, uid int64) ([]int64, error) {
	val, err := c.rdb.Get(ctx, followeesKey(uid)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, redis.Nil
	} else if err != nil {
		return nil, err
	}
	var res []int64
	err = json.Unmarshal([]byte(val), &res)
	return res, err
}

func (c *redisCache) SetFollowees(ctx context.Context, uid int64, res []int64) error {
	data, _ := json.Marshal(res)
	return c.rdb.Set(ctx, followeesKey(uid), data, c.listTTL).Err()
}

func (c *redisCache) GetFollowerCount(ctx context.Context, uid int64) (uint32, error) {
	val, err := c.rdb.HGet(ctx, followCountKey(uid), "follower_count").Int()
	if errors.Is(err, redis.Nil) {
		return 0, redis.Nil
	} else if err != nil {
		return 0, err
	}
	return uint32(val), nil
}

func (c *redisCache) GetFolloweeCount(ctx context.Context, uid int64) (uint32, error) {
	val, err := c.rdb.HGet(ctx, followCountKey(uid), "followee_count").Int()
	if errors.Is(err, redis.Nil) {
		return 0, redis.Nil
	} else if err != nil {
		return 0, err
	}
	return uint32(val), nil
}

func (c *redisCache) SetFollowersCount(ctx context.Context, uid int64, cnt uint32) error {
	return c.rdb.HSet(ctx, followCountKey(uid), "follower_count", cnt).Err()
}

func (c *redisCache) SetFolloweesCount(ctx context.Context, uid int64, cnt uint32) error {
	return c.rdb.HSet(ctx, followCountKey(uid), "followee_count", cnt).Err()
}
func (c *redisCache) UpdateRelation(ctx context.Context, info domain.RelationInfo, delta int) error {
	_, err := c.rdb.Eval(ctx, updateRelationLua, []string{followCountKey(info.Follower)}, delta, int(c.countTTL.Seconds())).Result()
	return err
}
