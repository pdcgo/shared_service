package auth_srv

import (
	"context"
	"encoding/base64"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"github.com/golang/protobuf/proto"
	"github.com/pdcgo/schema/services/access_iface/v1"
	"github.com/pdcgo/schema/services/user_iface/v1"
	"github.com/pdcgo/shared/db_models"
)

// CheckLogin implements user_ifaceconnect.AuthServiceHandler.
func (a *authServiceImpl) CheckLogin(
	ctx context.Context,
	req *connect.Request[user_iface.CheckLoginRequest],
) (*connect.Response[user_iface.CheckLoginResponse], error) {
	var err error

	identity := a.auth.AuthIdentityFromHeader(req.Header())
	agent := identity.Identity()
	err = identity.Err()
	if err != nil {
		return nil, err
	}

	db := a.db.WithContext(ctx)
	var usr user_iface.User
	err = db.Model(&user_iface.User{}).
		First(&usr, agent.IdentityID()).
		Error

	if err != nil {
		return nil, err
	}

	rawsource := req.Header().Get("X-Pdc-Source")
	// decoding source
	data, err := base64.StdEncoding.DecodeString(rawsource)

	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	source := &access_iface.RequestSource{}
	err = proto.Unmarshal(data, source)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	validator := protovalidate.GlobalValidator
	err = validator.Validate(source)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	var team db_models.Team
	err = db.Model(&db_models.Team{}).
		First(&team, source.TeamId).
		Error

	if err != nil {
		return nil, err
	}

	result := user_iface.CheckLoginResponse{
		User: &usr,
		Team: &user_iface.Team{
			Id:   uint64(team.ID),
			Name: team.Name,
		},
	}

	return connect.NewResponse(&result), nil

}
