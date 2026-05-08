//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"
	"github.com/pdcgo/san_collection/san_mcp"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/shared_service"
	"github.com/pdcgo/shared_service/services/user_service"
)

func InitializeApp() (*App, error) {
	wire.Build(
		http.NewServeMux,
		configs.NewProductionConfig,
		NewFirestoreClient,
		NewDatabase,
		NewRedisDatabase,
		NewCache,
		NewAuthorization,

		san_mcp.NewMcpSessionManager,
		custom_connect.NewDefaultInterceptor,
		custom_connect.NewRegisterReflect,
		user_service.NewRegister,
		shared_service.NewRegister,
		NewApp,
	)

	return &App{}, nil
}
