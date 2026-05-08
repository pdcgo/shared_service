package auth_srv

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pdcgo/san_collection/san_mcp"
	"github.com/pdcgo/schema/services/user_iface/v1"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"gorm.io/gorm"
)

type authServiceImpl struct {
	db             *gorm.DB
	auth           authorization_iface.Authorization
	sessionManager *san_mcp.McpSessionManager
	secret         string
}

// Logout implements user_ifaceconnect.AuthServiceHandler.
func (a *authServiceImpl) Logout(context.Context, *connect.Request[user_iface.LogoutRequest]) (*connect.Response[user_iface.LogoutResponse], error) {
	panic("unimplemented")
}

func NewAuthService(
	db *gorm.DB,
	auth authorization_iface.Authorization,
	sessionManager *san_mcp.McpSessionManager,
	secret string,
) *authServiceImpl {
	return &authServiceImpl{db, auth, sessionManager, secret}
}
