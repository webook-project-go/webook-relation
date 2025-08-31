//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/webook-project-go/webook-relation/grpc"
	"github.com/webook-project-go/webook-relation/ioc"
	"github.com/webook-project-go/webook-relation/repository"
	"github.com/webook-project-go/webook-relation/repository/cache"
	"github.com/webook-project-go/webook-relation/repository/dao"
	"github.com/webook-project-go/webook-relation/service"
)

var interactiveServiceProvider = wire.NewSet(
	service.New,
	repository.New,
	cache.New,
	dao.New,
)

var thirdPartyProvider = wire.NewSet(
	ioc.InitDatabase,
	ioc.InitRedis,
	ioc.InitEtcd,
)

func InitApp() *App {
	wire.Build(
		wire.Struct(new(App), "*"),
		thirdPartyProvider,
		grpc.New,
		interactiveServiceProvider,
		ioc.InitGrpcServer,
	)
	return new(App)
}
