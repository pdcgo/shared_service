package common

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/common/v1"
	"github.com/pdcgo/shared/db_models"
	"gorm.io/gorm"
)

type customerDataImpl struct {
	db *gorm.DB
}

// CustomerIDs implements commonconnect.CustomerDataServiceHandler.
func (c *customerDataImpl) CustomerIDs(
	ctx context.Context,
	req *connect.Request[common.CustomerIDsRequest]) (*connect.Response[common.CustomerIDsResponse], error) {
	var err error
	db := c.db.WithContext(ctx)
	pay := req.Msg
	result := common.CustomerIDsResponse{
		Data: map[uint64]*common.Customer{},
	}

	custs := []*common.Customer{}
	err = db.
		Model(&db_models.CustomerAddress{}).
		Where("id in ?", pay.Ids).
		Find(&custs).
		Error

	if err != nil {
		return nil, err
	}

	for _, d := range custs {
		result.Data[d.Id] = d
	}

	return connect.NewResponse(&result), nil

}

func NewCustomerDataService(db *gorm.DB) *customerDataImpl {
	return &customerDataImpl{
		db: db,
	}
}
