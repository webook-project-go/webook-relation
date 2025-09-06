package grpc

import (
	"context"
	"github.com/webook-project-go/webook-apis/gen/go/apis/relation/v1"
	"github.com/webook-project-go/webook-relation/domain"
	"github.com/webook-project-go/webook-relation/service"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Service struct {
	svc service.Service
	v1.UnimplementedRelationServiceServer
}

func New(svc service.Service) *Service {
	return &Service{svc: svc}
}
func (s *Service) Follow(ctx context.Context, info *v1.RelationInfo) (*emptypb.Empty, error) {
	now := time.Now().UnixMilli()
	err := s.svc.Follow(ctx, domain.RelationInfo{
		Status:   domain.StatusUnFollowing,
		Follower: info.GetFollower(),
		Followee: info.GetFollowee(),
		Ctime:    now,
		Utime:    now,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Service) UnFollow(ctx context.Context, info *v1.RelationInfo) (*emptypb.Empty, error) {
	err := s.svc.UnFollow(ctx, domain.RelationInfo{
		Follower: info.Follower,
		Followee: info.Followee,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Service) GetFollowers(ctx context.Context, req *v1.GetFollowersReq) (*v1.UserIDList, error) {
	res, err := s.svc.GetFollowers(ctx, req.GetUid(), req.GetLastID(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &v1.UserIDList{
		Uids: res,
	}, nil
}

func (s *Service) GetFollowees(ctx context.Context, req *v1.GetFolloweesReq) (*v1.UserIDList, error) {
	res, err := s.svc.GetFollowees(ctx, req.GetUid(), req.GetLastID(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &v1.UserIDList{Uids: res}, nil
}

func (s *Service) GetFollowerCount(ctx context.Context, req *v1.GetFollowerCountReq) (*v1.FollowerCountResp, error) {
	res, err := s.svc.GetFollowerCount(ctx, req.GetUid())
	if err != nil {
		return nil, err
	}
	return &v1.FollowerCountResp{Count: res}, nil
}

func (s *Service) GetFolloweeCount(ctx context.Context, req *v1.GetFolloweeCountReq) (*v1.FolloweeCountResp, error) {
	res, err := s.svc.GetFolloweeCount(ctx, req.GetUid())
	if err != nil {
		return nil, err
	}
	return &v1.FolloweeCountResp{Count: res}, nil
}
