package shared_service

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1/access_ifaceconnect"
	"github.com/pdcgo/schema/services/common/v1/commonconnect"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared_service/services/access_service"
	"github.com/pdcgo/shared_service/services/common"
	"github.com/pdcgo/shared_service/services/configuration"
	"github.com/pdcgo/shared_service/services/hello_service"
	"gorm.io/gorm"
)

type MigrationHandler func() error
type ServiceReflectNames []string
type RegisterHandler func() ServiceReflectNames

func NewRegister(
	mux *http.ServeMux,
	db *gorm.DB,
	cfg *configs.AppConfig,
	auth authorization_iface.Authorization,
	firestoreClient *firestore.Client,
	defaultInterceptor custom_connect.DefaultInterceptor,
) RegisterHandler {

	return func() ServiceReflectNames {
		grpcReflects := ServiceReflectNames{}
		path, handler := access_ifaceconnect.NewFrontendAccessServiceHandler(access_service.NewAccessService(db, auth), defaultInterceptor)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, access_ifaceconnect.FrontendAccessServiceName)

		path, handler = commonconnect.NewTeamServiceHandler(common.NewTeamService(db), defaultInterceptor)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, commonconnect.TeamServiceName)

		path, handler = commonconnect.NewShopServiceHandler(common.NewShopService(db), defaultInterceptor)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, commonconnect.ShopServiceName)

		path, handler = commonconnect.NewUserServiceHandler(common.NewUserService(db), defaultInterceptor)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, commonconnect.UserServiceName)

		path, handler = commonconnect.NewWarehouseServiceHandler(common.NewWarehouseService(db), defaultInterceptor)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, commonconnect.WarehouseServiceName)

		path, handler = commonconnect.NewCustomerDataServiceHandler(common.NewCustomerDataService(db), defaultInterceptor)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, commonconnect.CustomerDataServiceName)

		path, handler = commonconnect.NewShipmentServiceHandler(common.NewShipmentService(db), defaultInterceptor)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, commonconnect.ShipmentServiceName)

		// global configuration service
		path, handler = access_ifaceconnect.NewConfigurationServiceHandler(
			configuration.NewConfigurationService(auth, firestoreClient, cfg.GithubToken),
		)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, access_ifaceconnect.ConfigurationServiceName)

		// custom source

		path, handler = access_ifaceconnect.NewHelloServiceHandler(hello_service.NewHelloService(),
			defaultInterceptor,
			connect.WithInterceptors(
				&custom_connect.RequestSourceInterceptor{},
				&custom_connect.ScopeIntercept{},
			),
		)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, access_ifaceconnect.HelloServiceName)

		return grpcReflects
	}
}
