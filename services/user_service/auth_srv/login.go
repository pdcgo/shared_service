package auth_srv

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
	"github.com/pdcgo/schema/services/user_iface/v1"
	"github.com/pdcgo/shared/authorization"
	"github.com/pdcgo/shared/db_models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// Login implements user_ifaceconnect.AuthServiceHandler.
func (a *authServiceImpl) Login(
	ctx context.Context,
	req *connect.Request[user_iface.LoginRequest],
) (*connect.Response[user_iface.LoginResponse], error) {
	var err error

	pay := req.Msg
	username := pay.Username

	db := a.db.WithContext(ctx)

	var usr db_models.User
	err = db.Model(&db_models.User{}).
		Where("username = ?", username).
		First(&usr).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("anda bukan user kami")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(pay.Password))
	if err != nil {
		return nil, fmt.Errorf("username atau password salah")
	}

	// creating token
	now := time.Now()
	validUntil := now.Add(time.Hour * 24)
	identity := authorization.JwtIdentity{
		UserID:     usr.ID,
		SuperUser:  usr.IsSuperUser(),
		CreatedAt:  now.UnixMicro(),
		ValidUntil: validUntil.UnixMicro(),
	}

	token, err := identity.Serialize(a.secret)
	if err != nil {
		return nil, err
	}

	// checking team id jika payload ateam_id di set
	var teamId uint
	if pay.TeamId != 0 {
		err = db.Model(&db_models.UserTeam{}).
			Select("team_id").
			Where("user_id = ?", usr.ID).
			Where("team_id = ?", pay.TeamId).
			Find(&teamId).
			Error

		if err != nil {
			return nil, err
		}

		if teamId == 0 {
			return nil, fmt.Errorf("tidak punya akses ke team")
		}
	} else {
		err = db.Model(&db_models.UserTeam{}).
			Select("team_id").
			Where("user_id = ?", usr.ID).
			Limit(1).
			Find(&teamId).
			Error

		if err != nil {
			return nil, err
		}

		if teamId == 0 {
			return nil, fmt.Errorf("[default] tidak punya akses ke team")
		}
	}

	// getting team
	var team db_models.Team
	err = db.Model(&db_models.Team{}).
		Where("id = ?", teamId).
		First(&team).
		Error

	if err != nil {
		return nil, err
	}

	// creating x pdc source
	source := access_iface.RequestSource{
		TeamId:      uint64(teamId),
		RequestFrom: pay.From,
	}

	rawsource, err := proto.Marshal(&source)
	if err != nil {
		return nil, err
	}

	xPdcSourceToken := base64.StdEncoding.EncodeToString(rawsource)

	return &connect.Response[user_iface.LoginResponse]{
		Msg: &user_iface.LoginResponse{
			Token:      token,
			XPdcSource: xPdcSourceToken,
			User: &user_iface.User{
				Id:             uint64(usr.ID),
				Name:           usr.Name,
				Username:       usr.Username,
				ProfilePicture: usr.ProfilePicture,
			},
			Team: &user_iface.Team{
				Id:   uint64(team.ID),
				Name: team.Name,
			},
		},
	}, nil

}
