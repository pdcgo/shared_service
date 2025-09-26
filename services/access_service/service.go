package access_service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
	"github.com/pdcgo/shared/db_models"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"gorm.io/gorm"
)

type MenuKey string

const (
	StatMenu        MenuKey = "stat_menu"
	StatProductMenu MenuKey = "stat_product_menu"
	StatOrderMenu   MenuKey = "stat_order_menu"
)

type accessServiceImpl struct {
	db   *gorm.DB
	auth authorization_iface.Authorization
}

// SetupAccess implements access_ifaceconnect.FrontendAccessServiceHandler.
func (a *accessServiceImpl) SetupAccess(context.Context, *connect.Request[access_iface.SetupAccessRequest], *connect.ServerStream[access_iface.SetupAccessResponse]) error {
	panic("unimplemented")
}

// MenuAccess implements access_ifaceconnect.FrontendAccessServiceHandler.
func (a *accessServiceImpl) MenuAccess(
	ctx context.Context,
	req *connect.Request[access_iface.MenuAccessRequest],
) (*connect.Response[access_iface.MenuAccessResponse], error) {
	var err error

	pay := req.Msg
	db := a.db.WithContext(ctx)

	err = a.
		auth.
		AuthIdentityFromHeader(req.Header()).
		Err()

	if err != nil {
		return &connect.Response[access_iface.MenuAccessResponse]{}, err
	}

	var feature db_models.TeamFeature

	err = db.
		Model(&db_models.TeamFeature{}).
		Where("team_id = ?", pay.TeamId).
		Find(&feature).
		Error

	if err != nil {
		return &connect.Response[access_iface.MenuAccessResponse]{}, err
	}

	var policy access_iface.Policy
	if feature.ProductPriority {
		policy = access_iface.Policy_POLICY_ALLOW
	} else {
		policy = access_iface.Policy_POLICY_DENIED
	}

	result := access_iface.MenuAccessResponse{
		Data: map[string]*access_iface.AccessItem{
			string(StatMenu): {
				Policy: policy,
			},
			string(StatProductMenu): {
				Policy: policy,
			},
			string(StatOrderMenu): {
				Policy: policy,
			},
		},
	}

	return connect.NewResponse(&result), err
}

func NewAccessService(db *gorm.DB, auth authorization_iface.Authorization) *accessServiceImpl {
	return &accessServiceImpl{
		db:   db,
		auth: auth,
	}
}
