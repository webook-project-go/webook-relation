package main

import (
	"context"
	v1 "github.com/webook-project-go/webook-apis/gen/go/apis/relation/v1"
	_ "github.com/webook-project-go/webook-relation/config"
	"github.com/webook-project-go/webook-relation/ioc"
)

func main() {
	app := InitApp()
	shutdwon := ioc.InitOTEL()
	defer shutdwon(context.Background())
	v1.RegisterRelationServiceServer(app.Server, app.Service)
	err := app.Server.Serve()
	if err != nil {
		panic(err)
	}
}
