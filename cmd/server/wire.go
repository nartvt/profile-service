//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	coreConf "github.com/nartvt/go-core/conf"
	"github.com/nartvt/go-core/server"
	coreService "github.com/nartvt/go-core/service"
	"github.com/nartvt/profile-service/internal/biz"
	"github.com/nartvt/profile-service/internal/conf"
	data "github.com/nartvt/profile-service/internal/data"
	"github.com/nartvt/profile-service/internal/service"
)

// initApp init kratos application.
func initApp(*coreConf.Server, *conf.Data, log.Logger) (coreService.Service, func(), error) {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, service.ProviderSet, server.ProviderSet, initService))
}
