package common

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/common/v1"
	"github.com/pdcgo/shared/db_models"
	"gorm.io/gorm"
)

type warehouseServiceImpl struct {
	db *gorm.DB
}

type WarehouseList []*db_models.Warehouse

func (l WarehouseList) MapProto() map[uint64]*common.Warehouse {
	res := map[uint64]*common.Warehouse{}
	for _, item := range l {
		res[uint64(item.ID)] = &common.Warehouse{
			Id:   uint64(item.ID),
			Name: item.Name,
			Desc: item.Desc,
		}
	}
	return res
}

// PublicWarehouseIDs implements commonconnect.WarehouseServiceHandler.
func (w *warehouseServiceImpl) PublicWarehouseIDs(
	ctx context.Context,
	req *connect.Request[common.PublicWarehouseIDsRequest]) (*connect.Response[common.PublicWarehouseIDsResponse], error) {
	var err error

	db := w.db.WithContext(ctx)
	pay := req.Msg

	list := WarehouseList{}
	err = db.
		Model(&db_models.Warehouse{}).
		Where("id IN ?", pay.Ids).
		Find(&list).
		Error

	result := common.PublicWarehouseIDsResponse{
		Data: list.MapProto(),
	}

	return connect.NewResponse(&result), err

}

func NewWarehouseService(db *gorm.DB) *warehouseServiceImpl {
	return &warehouseServiceImpl{
		db: db,
	}
}
