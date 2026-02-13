package configuration

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
)

type configrationServiceImpl struct {
	ghToken string
	auth    authorization_iface.Authorization
	client  *firestore.Client
}

// ExtensionConfiguration implements access_ifaceconnect.ConfigurationServiceHandler.
func (c *configrationServiceImpl) ExtensionConfiguration(
	ctx context.Context,
	req *connect.Request[access_iface.ExtensionConfigurationRequest],
) (*connect.Response[access_iface.ExtensionConfigurationResponse], error) {
	pay := req.Msg
	mode := access_iface.Mode_name[int32(pay.Mode)]
	version := pay.Version
	ref := fmt.Sprintf("%s/extension_configuration/%s", mode, version)
	doc := c.client.Collection("configuration").Doc(ref)
	snap, err := doc.Get(ctx)
	if err != nil {
		return nil, err
	}

	result := access_iface.ExtensionConfigurationResponse{
		Data: &access_iface.ExtensionConfiguration{},
	}

	err = snap.DataTo(result.Data)
	return connect.NewResponse(&result), err
}

// ExtensionConfigurationReplace implements access_ifaceconnect.ConfigurationServiceHandler.
func (c *configrationServiceImpl) ExtensionConfigurationReplace(
	ctx context.Context,
	req *connect.Request[access_iface.ExtensionConfigurationReplaceRequest],
) (*connect.Response[access_iface.ExtensionConfigurationReplaceResponse], error) {
	var err error

	pay := req.Msg
	conf := pay.Data

	identity := c.auth.AuthIdentityFromToken(pay.Token)
	agent := identity.Identity()
	err = identity.Err()
	if err != nil {
		return nil, err
	}

	if !agent.IsSuperUser() {
		return nil, fmt.Errorf("anda siapa ya ?")
	}

	// creating configuration on firebase

	// building ref document
	mode := access_iface.Mode_name[int32(conf.Mode)]
	version := conf.Version
	ref := fmt.Sprintf("%s/extension_configuration/%s", mode, version)
	doc := c.client.Collection("configuration").Doc(ref)
	_, err = doc.Set(ctx, conf)
	if err != nil {
		return nil, err
	}

	return &connect.Response[access_iface.ExtensionConfigurationReplaceResponse]{
		Msg: &access_iface.ExtensionConfigurationReplaceResponse{},
	}, nil

}

func NewConfigurationService(
	auth authorization_iface.Authorization,
	client *firestore.Client,
	ghToken string,
) *configrationServiceImpl {
	return &configrationServiceImpl{ghToken, auth, client}
}
