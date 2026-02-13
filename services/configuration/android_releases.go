package configuration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
)

// AndroidReleases implements [access_ifaceconnect.ConfigurationServiceHandler].
func (c *configrationServiceImpl) AndroidReleases(
	ctx context.Context,
	req *connect.Request[access_iface.AndroidReleasesRequest]) (*connect.Response[access_iface.AndroidReleasesResponse], error) {
	var err error

	hreq, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/PDC-Repository/warehouse_android_flutter/releases", nil)
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

	res := access_iface.AndroidReleasesResponse{
		Releases: []*access_iface.Release{},
	}
	err = json.Unmarshal(raw, &res.Releases)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&res), err
}
