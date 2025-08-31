package dao

import (
	"context"
	"errors"
	"github.com/webook-project-go/webook-relation/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Dao interface {
	UpsertFollowInfo(ctx context.Context, entity RelationInfo) error
	MarkUnFollow(ctx context.Context, entity RelationInfo) error
	FindFollowers(ctx context.Context, uid int64) ([]int64, error)
	FindFollowees(ctx context.Context, uid int64) ([]int64, error)
	GetFollowerCount(ctx context.Context, uid int64) (uint32, error)
	GetFolloweeCount(ctx context.Context, uid int64) (uint32, error)
}
type dao struct {
	db *gorm.DB
}

func New(db *gorm.DB) Dao {
	return &dao{db: db}
}

func (d *dao) UpsertFollowInfo(ctx context.Context, entity RelationInfo) error {
	return d.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "follower"}, {Name: "followee"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status": entity.Status,
				"utime":  entity.Utime,
			}),
		},
	).Create(&entity).Error
}

func (d *dao) MarkUnFollow(ctx context.Context, entity RelationInfo) error {
	result := d.db.WithContext(ctx).Model(&RelationInfo{}).
		Where("follower = ? AND followee = ? AND status = ?", entity.Follower, entity.Followee, domain.StatusFollowing).
		Update("status", domain.StatusUnFollowing)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no such relation")
	}
	return nil
}

func (d *dao) FindFollowers(ctx context.Context, uid int64) ([]int64, error) {
	var res []RelationInfo
	err := d.db.WithContext(ctx).
		Where("followee = ? AND status = ?", uid, domain.StatusFollowing).
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(res))
	for _, r := range res {
		ids = append(ids, r.Follower)
	}
	return ids, nil
}

func (d *dao) FindFollowees(ctx context.Context, uid int64) ([]int64, error) {
	var res []RelationInfo
	err := d.db.WithContext(ctx).
		Where("follower = ? AND status = ?", uid, domain.StatusFollowing).
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(res))
	for _, r := range res {
		ids = append(ids, r.Followee)
	}
	return ids, nil
}

func (d *dao) GetFollowerCount(ctx context.Context, uid int64) (uint32, error) {
	var info FollowCount
	err := d.db.WithContext(ctx).
		Where("uid = ?", uid).Find(&info).Error
	if err != nil {
		return 0, err
	}
	return info.FollowerCount, nil
}

func (d *dao) GetFolloweeCount(ctx context.Context, uid int64) (uint32, error) {
	var info FollowCount
	err := d.db.WithContext(ctx).
		Where("uid = ?", uid).Find(&info).Error
	if err != nil {
		return 0, err
	}
	return info.FolloweeCount, nil
}
