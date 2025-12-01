package shared_service

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1/access_ifaceconnect"
	"github.com/pdcgo/schema/services/common/v1/commonconnect"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared_service/services/access_service"
	"github.com/pdcgo/shared_service/services/common"
	"github.com/pdcgo/shared_service/services/hello_service"
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

		path, handler := access_ifaceconnect.NewFrontendAccessServiceHandler(access_service.NewAccessService(db, auth), defaultInterceptor)
		mux.Handle(path, handler)
		path, handler = commonconnect.NewTeamServiceHandler(common.NewTeamService(db), defaultInterceptor)
		mux.Handle(path, handler)
		path, handler = commonconnect.NewShopServiceHandler(common.NewShopService(db), defaultInterceptor)
		mux.Handle(path, handler)
		path, handler = commonconnect.NewUserServiceHandler(common.NewUserService(db), defaultInterceptor)
		mux.Handle(path, handler)
		path, handler = commonconnect.NewWarehouseServiceHandler(common.NewWarehouseService(db), defaultInterceptor)
		mux.Handle(path, handler)

		path, handler = commonconnect.NewCustomerDataServiceHandler(common.NewCustomerDataService(db), defaultInterceptor)
		mux.Handle(path, handler)

		// custom source

		path, handler = access_ifaceconnect.NewHelloServiceHandler(hello_service.NewHelloService(),
			defaultInterceptor,
			connect.WithInterceptors(
				&custom_connect.RequestSourceInterceptor{},
				&custom_connect.ScopeIntercept{},
			),
		)
		mux.Handle(path, handler)

	}
}
