package configuration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
)

// AndroidCheckLatestVersion implements [access_ifaceconnect.ConfigurationServiceHandler].
func (c *configrationServiceImpl) AndroidCheckLatestVersion(
	ctx context.Context,
	req *connect.Request[access_iface.AndroidCheckLatestVersionRequest],
) (*connect.Response[access_iface.AndroidCheckLatestVersionResponse], error) {
	hreq, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/PDC-Repository/warehouse_android_flutter/releases/latest", nil)
	if err != nil {
		return nil, err
	}

	hreq.Header.Add("Authorization", "Bearer "+c.ghToken)
	hres, err := http.DefaultClient.Do(hreq)
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(hres.Body)
	if err != nil {
		return nil, err
	}

	res := access_iface.AndroidCheckLatestVersionResponse{
		Release: &access_iface.Release{},
	}

	// log.Println(string(raw))

	err = json.Unmarshal(raw, res.Release)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&res), err
}
