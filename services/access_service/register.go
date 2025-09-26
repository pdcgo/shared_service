package access_service

import (
	"net/http"

	"github.com/pdcgo/schema/services/access_iface/v1/access_ifaceconnect"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"gorm.io/gorm"
)

type MigrationHandler func() error

type RegisterHandler func()

func NewRegister(
	mux *http.ServeMux,
	db *gorm.DB,
	auth authorization_iface.Authorization,
	defaultInterceptor custom_connect.DefaultInterceptor,
) RegisterHandler {

	return func() {

		path, handler := access_ifaceconnect.NewFrontendAccessServiceHandler(NewAccessService(db, auth), defaultInterceptor)
		mux.Handle(path, handler)

	}
}
