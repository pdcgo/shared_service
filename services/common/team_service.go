package common

import (
	"context"
	"math"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/common/v1"
	"github.com/pdcgo/shared/db_models"
	"gorm.io/gorm"
)

type teamServiceImpl struct {
	db *gorm.DB
}

// PublicTeamList implements commonconnect.TeamServiceHandler.
func (t *teamServiceImpl) PublicTeamList(
	ctx context.Context,
	req *connect.Request[common.PublicTeamListRequest],
) (*connect.Response[common.PublicTeamListResponse], error) {
	var err error

	pay := req.Msg
	db := t.db.WithContext(ctx)
	result := common.PublicTeamListResponse{
		Datas:    []*common.Team{},
		PageInfo: &common.PageInfo{},
	}

	page := pay.Page
	offset := page.Page*page.Limit - page.Limit

	tx := db.Model(&db_models.Team{}).Where("deleted = ?", false)
	if pay.Q != "" {
		like := "%" + pay.Q + "%"
		tx = tx.Where("name ILIKE ? OR team_code ILIKE ?", like, like)
	}
	switch pay.TeamType {
	case common.TeamType_TEAM_TYPE_WAREHOUSE:
		tx = tx.Where("type = ?", db_models.WarehouseTeamType)
	case common.TeamType_TEAM_TYPE_SELLING:
		tx = tx.Where("type = ?", db_models.SellingTeamType)
	case common.TeamType_TEAM_TYPE_ADMIN:
		tx = tx.Where("type = ?", db_models.AdminTeamType)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return &connect.Response[common.PublicTeamListResponse]{}, err
	}

	err = tx.
		Select([]string{
			"id",
			"name",
			"team_code",
			"type",
		}).
		Limit(int(page.Limit)).
		Offset(int(offset)).
		Find(&result.Datas).Error

	if err != nil {
		return nil, err
	}

	pageCount := math.Ceil(float64(total) / float64(page.Limit))
	result.PageInfo.CurrentPage = page.Page
	result.PageInfo.TotalPage = int64(pageCount)
	result.PageInfo.TotalItems = total

	return connect.NewResponse(&result), nil
}

// PublicTeamIDs implements commonconnect.TeamServiceHandler.
func (t *teamServiceImpl) PublicTeamIDs(
	ctx context.Context,
	req *connect.Request[common.PublicTeamIDsRequest],
) (*connect.Response[common.PublicTeamIDsResponse], error) {
	var err error
	result := common.PublicTeamIDsResponse{
		Data: map[uint64]*common.Team{},
	}

	db := t.db.WithContext(ctx)
	pay := req.Msg

	teams := []*db_models.Team{}

	err = db.
		Model(db_models.Team{}).
		Preload("TeamInfo").
		Where("id in ?", pay.Ids).
		Find(&teams).
		Error

	for _, d := range teams {
		info := &common.TeamInfo{}
		if d.TeamInfo != nil {
			info = &common.TeamInfo{
				TeamId:        uint64(d.ID),
				ContactNumber: d.TeamInfo.ContactNumber,
			}

			if d.TeamInfo.ReturnWarehouseID != nil {
				info.ReturnWarehouseId = uint64(*d.TeamInfo.ReturnWarehouseID)
			}

			if d.TeamInfo.ReturnUserID != nil {
				info.ReturnUserId = uint64(*d.TeamInfo.ReturnUserID)
			}

		}

		result.Data[uint64(d.ID)] = &common.Team{
			Id:       uint64(d.ID),
			Name:     d.Name,
			TeamCode: string(d.TeamCode),
			Type:     string(d.Type),
			Info:     info,
		}
	}

	return connect.NewResponse(&result), err
}

func NewTeamService(db *gorm.DB) *teamServiceImpl {
	return &teamServiceImpl{
		db: db,
	}
}
