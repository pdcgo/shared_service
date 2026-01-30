package common

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/common/v1"
	"github.com/pdcgo/shared/db_models"
	"gorm.io/gorm"
)

type shipmentServiceImpl struct {
	db *gorm.DB
}

// PublicShipmentList implements commonconnect.ShipmentServiceHandler.
func (s *shipmentServiceImpl) PublicShipmentList(
	ctx context.Context,
	req *connect.Request[common.PublicShipmentListRequest],
) (*connect.Response[common.PublicShipmentListResponse], error) {
	var err error

	db := s.db.WithContext(ctx)
	result := &common.PublicShipmentListResponse{
		Data: []*common.Shipment{},
	}

	err = db.
		Model(&db_models.Shipping{}).
		Find(&result.Data).
		Error

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(result), nil

}

// PublicShipmentIDs implements commonconnect.ShipmentServiceHandler.
func (s *shipmentServiceImpl) PublicShipmentIDs(
	ctx context.Context,
	req *connect.Request[common.PublicShipmentIDsRequest],
) (*connect.Response[common.PublicShipmentIDsResponse], error) {
	var err error

	db := s.db.WithContext(ctx)
	pay := req.Msg

	result := common.PublicShipmentIDsResponse{
		Data: map[uint64]*common.Shipment{},
	}

	list := []*common.Shipment{}

	err = db.
		Model(&db_models.Shipping{}).
		Where("id in ?", pay.Ids).
		Find(&list).
		Error

	if err != nil {
		return nil, err
	}

	for _, item := range list {
		result.Data[uint64(item.Id)] = item
	}

	return connect.NewResponse(&result), nil

}

func NewShipmentService(db *gorm.DB) *shipmentServiceImpl {
	return &shipmentServiceImpl{
		db: db,
	}
}
