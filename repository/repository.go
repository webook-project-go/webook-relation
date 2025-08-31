package repository

import (
	"context"
	"github.com/webook-project-go/webook-relation/domain"
	"github.com/webook-project-go/webook-relation/repository/cache"
	"github.com/webook-project-go/webook-relation/repository/dao"
)

type Repository interface {
	MarkFollow(ctx context.Context, info domain.RelationInfo) error
	MarkUnFollow(ctx context.Context, info domain.RelationInfo) error
	FindFollowers(ctx context.Context, uid int64) ([]int64, error)
	FindFollowees(ctx context.Context, uid int64) ([]int64, error)
	GetFollowerCount(ctx context.Context, uid int64) (uint32, error)
	GetFolloweeCount(ctx context.Context, uid int64) (uint32, error)
}

type repository struct {
	cache cache.Cache
	d     dao.Dao
}

func New(cache cache.Cache, d dao.Dao) Repository {
	return &repository{
		cache: cache,
		d:     d,
	}
}
func RelationToEntity(info domain.RelationInfo) dao.RelationInfo {
	return dao.RelationInfo{
		Follower: info.Follower,
		Followee: info.Followee,
		Status:   info.Status,
	}
}
func RelationToDomain(info dao.RelationInfo) domain.RelationInfo {
	return domain.RelationInfo{
		ID:       info.ID,
		Follower: info.Follower,
		Followee: info.Followee,
		Status:   info.Status,
		Ctime:    info.Ctime,
		Utime:    info.Utime,
	}
}
func (r *repository) MarkFollow(ctx context.Context, info domain.RelationInfo) error {
	err := r.d.UpsertFollowInfo(ctx, RelationToEntity(info))
	if err != nil {
		return err
	}
	_ = r.cache.UpdateRelation(ctx, info, 1)
	return nil
}

func (r *repository) MarkUnFollow(ctx context.Context, info domain.RelationInfo) error {
	err := r.d.MarkUnFollow(ctx, RelationToEntity(info))
	if err != nil {
		return err
	}
	_ = r.cache.UpdateRelation(ctx, info, -1)
	return nil
}

func (r *repository) FindFollowers(ctx context.Context, uid int64) ([]int64, error) {
	res, err := r.cache.FindFollowers(ctx, uid)
	if err == nil && len(res) > 0 {
		return res, nil
	}
	res, err = r.d.FindFollowers(ctx, uid)
	if err != nil {
		return nil, err
	}
	_ = r.cache.SetFollowers(ctx, uid, res)
	return res, nil
}

func (r *repository) FindFollowees(ctx context.Context, uid int64) ([]int64, error) {
	res, err := r.cache.FindFollowees(ctx, uid)
	if err == nil && len(res) > 0 {
		return res, nil
	}
	res, err = r.d.FindFollowees(ctx, uid)
	if err != nil {
		return nil, err
	}
	_ = r.cache.SetFollowees(ctx, uid, res)
	return res, nil
}

func (r *repository) GetFollowerCount(ctx context.Context, uid int64) (uint32, error) {
	cnt, err := r.cache.GetFollowerCount(ctx, uid)
	if err == nil && cnt != 0 {
		return cnt, nil
	}
	cnt, err = r.d.GetFollowerCount(ctx, uid)
	if err != nil {
		return 0, err
	}
	_ = r.cache.SetFollowersCount(ctx, uid, cnt)
	return cnt, nil
}

func (r *repository) GetFolloweeCount(ctx context.Context, uid int64) (uint32, error) {
	cnt, err := r.cache.GetFolloweeCount(ctx, uid)
	if err == nil && cnt != 0 {
		return cnt, nil
	}
	cnt, err = r.d.GetFolloweeCount(ctx, uid)
	if err != nil {
		return 0, err
	}
	_ = r.cache.SetFolloweesCount(ctx, uid, cnt)
	return cnt, nil
}
