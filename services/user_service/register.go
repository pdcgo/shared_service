package user_service

import (
	"net/http"

	"github.com/pdcgo/schema/services/user_iface/v1/user_ifaceconnect"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared_service/services/user_service/auth_srv"

	"gorm.io/gorm"
)

type ServiceReflectNames []string
type RegisterHandler func() ServiceReflectNames

func NewRegister(
	db *gorm.DB,
	cfg *configs.AppConfig,
	auth authorization_iface.Authorization,
	mux *http.ServeMux,
	defaultInterceptor custom_connect.DefaultInterceptor,
	// cache ware_cache.Cache,
	// dispather report.ReportDispatcher,
) RegisterHandler {
	return func() ServiceReflectNames {
		grpcReflects := ServiceReflectNames{}

		path, handler := user_ifaceconnect.NewAuthServiceHandler(auth_srv.NewAuthService(db, auth, cfg.JwtSecret))
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, user_ifaceconnect.AuthServiceName)

		return grpcReflects
	}
}
