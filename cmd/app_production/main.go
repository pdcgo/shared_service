package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/pdcgo/shared/authorization"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/shared/db_connect"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared/pkg/cloud_logging"
	"github.com/pdcgo/shared/pkg/ware_cache"
	"github.com/pdcgo/shared_service"
	"github.com/pdcgo/shared_service/services/user_service"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"gorm.io/gorm"
)

func NewFirestoreClient() (*firestore.Client, error) {
	return firestore.NewClient(context.Background(), os.Getenv("GOOGLE_CLOUD_PROJECT"))
}

func NewDatabase(cfg *configs.AppConfig) (*gorm.DB, error) {
	return db_connect.NewProductionDatabase("shared_service", &cfg.Database)
}

func NewCache() (ware_cache.Cache, error) {
	return ware_cache.NewBadgerCache("/tmp/cache")
}

func NewAuthorization(
	cfg *configs.AppConfig,
	db *gorm.DB,
	cache ware_cache.Cache,
) authorization_iface.Authorization {
	return authorization.NewAuthorization(cache, db, cfg.JwtSecret)
}

type App struct {
	Run func() error
}

func NewApp(
	mux *http.ServeMux,
	accessRegister shared_service.RegisterHandler,
	userServiceRegister user_service.RegisterHandler,
	reflectRegister custom_connect.RegisterReflectFunc,
) *App {
	return &App{
		Run: func() error {
			cancel, err := custom_connect.InitTracer("shared-service")
			if err != nil {
				return err
			}

			defer cancel(context.Background())

			var grpcReflectNames []string

			grpcReflectNames = append(grpcReflectNames, accessRegister()...)
			grpcReflectNames = append(grpcReflectNames, userServiceRegister()...)

			reflectRegister(grpcReflectNames)

			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}

			host := os.Getenv("HOST")
			listen := fmt.Sprintf("%s:%s", host, port)
			log.Println("listening on", listen)

			return http.ListenAndServe(
				listen,
				// Use h2c so we can serve HTTP/2 without TLS.
				h2c.NewHandler(
					custom_connect.WithCORS(mux),
					&http2.Server{}),
			)
		},
	}
}

func main() {
	cloud_logging.SetCloudLoggingDefault()
	app, err := InitializeApp()
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
