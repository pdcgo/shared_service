//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/shared_service"
	"github.com/pdcgo/user_service"
)

func InitializeApp() (*App, error) {
	wire.Build(
		http.NewServeMux,
		configs.NewProductionConfig,
		NewDatabase,
		NewCache,
		NewAuthorization,
		custom_connect.NewDefaultInterceptor,
		custom_connect.NewRegisterReflect,
		user_service.NewRegister,
		shared_service.NewRegister,
		NewApp,
	)

	return &App{}, nil
}
