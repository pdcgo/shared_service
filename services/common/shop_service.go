package common

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/common/v1"
	"github.com/pdcgo/shared/db_models"
	"gorm.io/gorm"
)

type shopServiceImpl struct {
	db *gorm.DB
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
