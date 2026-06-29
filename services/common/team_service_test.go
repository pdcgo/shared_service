package common

import (
	"strings"
	"testing"

	"connectrpc.com/connect"
	common "github.com/pdcgo/schema/services/common/v1"
	"github.com/pdcgo/shared/db_models"
	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/pdcgo/shared/pkg/moretest/moretest_mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestPublicTeamListTeamTypeFilter covers the team_type filter added to PublicTeamList:
// UNSPECIFIED returns all types, a specific type returns only that type, and the filter
// composes with the q search.
func TestPublicTeamListTeamTypeFilter(t *testing.T) {
	var scenario moretest_mock.DbScenario
	moretest.Suite(t, "public team list team_type filter",
		moretest.SetupListFunc{moretest_mock.MockPostgresDatabase(&scenario)},
		func(t *testing.T) {
			scenario(t, func(db *gorm.DB) {
				assert.NoError(t, db.AutoMigrate(&db_models.Team{}))

				seed := func(name string, tt db_models.TeamType) {
					code := strings.ToUpper(strings.ReplaceAll(name, " ", ""))
					assert.NoError(t, db.Create(&db_models.Team{
						Name:     name,
						Type:     tt,
						TeamCode: db_models.TeamCode(code),
					}).Error)
				}
				seed("Ware One", db_models.WarehouseTeamType)
				seed("Ware Two", db_models.WarehouseTeamType)
				seed("Sell One", db_models.SellingTeamType)
				seed("Admin One", db_models.AdminTeamType)

				svc := NewTeamService(db)
				list := func(tt common.TeamType, q string) []*common.Team {
					res, err := svc.PublicTeamList(t.Context(), connect.NewRequest(&common.PublicTeamListRequest{
						Q:        q,
						TeamType: tt,
						Page:     &common.PageFilter{Page: 1, Limit: 100},
					}))
					assert.NoError(t, err)
					return res.Msg.Datas
				}

				t.Run("unspecified returns all types", func(t *testing.T) {
					assert.Len(t, list(common.TeamType_TEAM_TYPE_UNSPECIFIED, ""), 4)
				})

				t.Run("warehouse returns only warehouse teams", func(t *testing.T) {
					got := list(common.TeamType_TEAM_TYPE_WAREHOUSE, "")
					assert.Len(t, got, 2)
					for _, tm := range got {
						assert.Equal(t, string(db_models.WarehouseTeamType), tm.Type)
					}
				})

				t.Run("selling returns only selling teams", func(t *testing.T) {
					assert.Len(t, list(common.TeamType_TEAM_TYPE_SELLING, ""), 1)
				})

				t.Run("admin returns only admin teams", func(t *testing.T) {
					assert.Len(t, list(common.TeamType_TEAM_TYPE_ADMIN, ""), 1)
				})

				t.Run("q composes with team_type", func(t *testing.T) {
					// "One" matches Ware One, Sell One, Admin One; warehouse narrows to just one.
					got := list(common.TeamType_TEAM_TYPE_WAREHOUSE, "One")
					assert.Len(t, got, 1)
					assert.Equal(t, "Ware One", got[0].Name)
				})
			})
		},
	)
}
