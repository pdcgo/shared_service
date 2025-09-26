//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/shared_service/services/access_service"
)

func InitializeApp() (*App, error) {
	wire.Build(
		http.NewServeMux,
		configs.NewProductionConfig,
		NewDatabase,
		NewCache,
		NewAuthorization,
		custom_connect.NewDefaultInterceptor,
		access_service.NewRegister,
		NewApp,
	)

	return &App{}, nil
}
