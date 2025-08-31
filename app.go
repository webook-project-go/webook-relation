package main

import (
	"github.com/webook-project-go/webook-pkgs/grpcx"
	"github.com/webook-project-go/webook-relation/grpc"
)

type App struct {
	Server  *grpcx.GrpcxServer
	Service *grpc.Service
}
