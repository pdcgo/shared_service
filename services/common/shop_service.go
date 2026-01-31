package common

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/common/v1"
	"github.com/pdcgo/shared/db_models"
	"gorm.io/gorm"
)

type shopServiceImpl struct {
	db *gorm.DB
}

// PublicShopList implements commonconnect.ShopServiceHandler.
func (s *shopServiceImpl) PublicShopList(ctx context.Context, req *connect.Request[common.PublicShopListRequest]) (*connect.Response[common.PublicShopListResponse], error) {
	var err error
	db := s.db.WithContext(ctx)
	pay := req.Msg

	result := common.PublicShopListResponse{
		Data: []*common.Shop{},
	}

	query := db.
		Model(&db_models.Marketplace{})

	if pay.TeamId != 0 {
		query = query.Where("team_id = ?", pay.TeamId)
	}

	if pay.MarketplaceType != common.MarketplaceType_MARKETPLACE_TYPE_UNSPECIFIED {
		var mptype db_models.MarketplaceType

		switch pay.MarketplaceType {
		case common.MarketplaceType_MARKETPLACE_TYPE_CUSTOM:
			mptype = db_models.MpCustom
		case common.MarketplaceType_MARKETPLACE_TYPE_LAZADA:
			mptype = db_models.MpLazada
		case common.MarketplaceType_MARKETPLACE_TYPE_MENGANTAR:
			mptype = db_models.MpMengantar
		case common.MarketplaceType_MARKETPLACE_TYPE_SHOPEE:
			mptype = db_models.MpShopee
		case common.MarketplaceType_MARKETPLACE_TYPE_TOKOPEDIA:
			mptype = db_models.MpTokopedia
		case common.MarketplaceType_MARKETPLACE_TYPE_TIKTOK:
			mptype = db_models.MpTiktok
		default:
			return &connect.Response[common.PublicShopListResponse]{}, fmt.Errorf("%s not supported", pay.MarketplaceType)
		}

		query = query.
			Where("mp_type = ?", mptype)
	}

	if pay.Q != "" {
		q := "%" + strings.ToLower(pay.Q) + "%"
		query = query.
			Where("lower(mp_name) like ? or lower(mp_username) like ?", q, q)
	}

	if pay.UserId != 0 {
		query = query.Where("user_id = ?", pay.UserId)
	}

	query = query.
		Select([]string{
			"id",
			"team_id",
			"marketplaces.mp_name as shop_name",
			"marketplaces.mp_username as shop_username",
			`
			case 
				when marketplaces.mp_type = 'tokopedia' then 2
				when marketplaces.mp_type = 'shopee' then 3
				when marketplaces.mp_type = 'lazada' then 5
				when marketplaces.mp_type = 'mengantar' then 6
				when marketplaces.mp_type = 'tiktok' then 4
				when marketplaces.mp_type = 'custom' then 1
				else 0

			end as marketplace_type
			`,
			"uri",
			//   string shop_name = 3;
			//   string shop_username = 4;
			//   MarketplaceType marketplace_type = 5;
			//   string uri = 6;
		})

	err = query.
		Limit(int(pay.Limit)).
		Find(&result.Data).
		Error

	if err != nil {
		return &connect.Response[common.PublicShopListResponse]{}, err
	}

	return connect.NewResponse(&result), nil
}

// PublicShopIDs implements commonconnect.ShopServiceHandler.
func (s *shopServiceImpl) PublicShopIDs(
	ctx context.Context,
	req *connect.Request[common.PublicShopIDsRequest],
) (*connect.Response[common.PublicShopIDsResponse], error) {
	var err error

	db := s.db.WithContext(ctx)
	pay := req.Msg

	result := common.PublicShopIDsResponse{
		Data: map[uint64]*common.Shop{},
	}

	shops := []*common.Shop{}

	// MpShopee    MarketplaceType = "shopee"
	// MpTokopedia MarketplaceType = "tokopedia"
	// MpTiktok    MarketplaceType = "tiktok"
	// MpMengantar MarketplaceType = "mengantar"
	// MpCustom    MarketplaceType = "custom"
	// MpLazada    MarketplaceType = "lazada"

	err = db.
		Model(&db_models.Marketplace{}).
		Select([]string{
			"id",
			"team_id",
			"marketplaces.mp_name as shop_name",
			"marketplaces.mp_username as shop_username",
			`
			case 
				when marketplaces.mp_type = 'tokopedia' then 2
				when marketplaces.mp_type = 'shopee' then 3
				when marketplaces.mp_type = 'lazada' then 5
				when marketplaces.mp_type = 'mengantar' then 6
				when marketplaces.mp_type = 'tiktok' then 4
				when marketplaces.mp_type = 'custom' then 1
				else 0

			end as marketplace_type
			`,
			"uri",
			//   string shop_name = 3;
			//   string shop_username = 4;
			//   MarketplaceType marketplace_type = 5;
			//   string uri = 6;
		}).
		Where("id in ?", pay.Ids).
		Find(&shops).
		Error

	if err != nil {
		return connect.NewResponse(&result), err
	}

	for _, d := range shops {
		result.Data[d.Id] = d
	}

	return connect.NewResponse(&result), nil

}

func NewShopService(db *gorm.DB) *shopServiceImpl {
	return &shopServiceImpl{
		db: db,
	}
}
