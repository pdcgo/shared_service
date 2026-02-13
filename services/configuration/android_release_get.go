package configuration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
)

// AndroidReleaseGet implements [access_ifaceconnect.ConfigurationServiceHandler].
func (c *configrationServiceImpl) AndroidReleaseGet(
	ctx context.Context,
	req *connect.Request[access_iface.AndroidReleaseGetRequest],
) (*connect.Response[access_iface.AndroidReleaseGetResponse], error) {

	var err error
	var uri string

	pay := req.Msg

	switch pay.By.(type) {
	case *access_iface.AndroidReleaseGetRequest_ReleaseId:
		releaseId := pay.GetReleaseId()
		uri = fmt.Sprintf("https://api.github.com/repos/PDC-Repository/warehouse_android_flutter/releases/%d", releaseId)
	case *access_iface.AndroidReleaseGetRequest_Tag:
		tag := pay.GetTag()
		uri = fmt.Sprintf("https://api.github.com/repos/PDC-Repository/warehouse_android_flutter/releases/tags/%s", tag)
	}

	hreq, err := http.NewRequest(http.MethodGet, uri, nil)
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

	res := access_iface.AndroidReleaseGetResponse{
		Release: &access_iface.Release{},
	}

	// log.Println(string(raw))

	err = json.Unmarshal(raw, res.Release)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&res), err
}
