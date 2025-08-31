package service

import (
	"context"
	"errors"
	"github.com/webook-project-go/webook-relation/domain"
	"github.com/webook-project-go/webook-relation/repository"
)

type Service interface {
	Follow(ctx context.Context, info domain.RelationInfo) error
	UnFollow(ctx context.Context, info domain.RelationInfo) error

	GetFollowers(ctx context.Context, uid int64) ([]int64, error)
	GetFollowees(ctx context.Context, uid int64) ([]int64, error)

	GetFollowerCount(ctx context.Context, uid int64) (uint32, error)
	GetFolloweeCount(ctx context.Context, uid int64) (uint32, error)
}

func New(repo repository.Repository) Service {
	return &service{repo: repo}
}

type service struct {
	repo repository.Repository
}

func (s *service) Follow(ctx context.Context, info domain.RelationInfo) error {
	return s.repo.MarkFollow(ctx, info)
}

func (s *service) UnFollow(ctx context.Context, info domain.RelationInfo) error {
	return s.repo.MarkUnFollow(ctx, info)
}

func (s *service) GetFollowers(ctx context.Context, uid int64) ([]int64, error) {
	if uid <= 0 {
		return nil, errors.New("illegal uid")
	}
	return s.repo.FindFollowers(ctx, uid)
}

func (s *service) GetFollowees(ctx context.Context, uid int64) ([]int64, error) {
	if uid <= 0 {
		return nil, errors.New("illegal uid")
	}
	return s.repo.FindFollowees(ctx, uid)
}

func (s *service) GetFollowerCount(ctx context.Context, uid int64) (uint32, error) {
	if uid <= 0 {
		return 0, errors.New("illegal uid")
	}
	return s.repo.GetFollowerCount(ctx, uid)
}

func (s *service) GetFolloweeCount(ctx context.Context, uid int64) (uint32, error) {
	if uid <= 0 {
		return 0, errors.New("illegal uid")
	}
	return s.repo.GetFolloweeCount(ctx, uid)
}
