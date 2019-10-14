package hydra

import (
	"encoding/json"
	"github.com/labbsr0x/goh/gohclient"
	"github.com/labbsr0x/goh/gohtypes"
	"net/http"
	"net/url"
	"path"
)

func get(client *gohclient.Default, flow, challenge string) map[string]interface{} {
	p := path.Join(client.BaseURL.Path, "/oauth2/auth/requests/", flow) + "?challenge=" + url.QueryEscape(challenge)
	return treatResponse(client.Get(p))
}

func put(client *gohclient.Default, flow, challenge, action string, data []byte) map[string]interface{} {
	p := path.Join(client.BaseURL.Path, "/oauth2/auth/requests/", flow, action) + "?challenge=" + url.QueryEscape(challenge)
	return treatResponse(client.Put(p, data))
}

func treatResponse(resp *http.Response, data []byte, err error) map[string]interface{} {
	if err == nil {
		if resp.StatusCode >= 200 && resp.StatusCode <= 302 {
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err == nil {
				return result
			}
			panic(gohtypes.Error{Code: 500, Err: err, Message: "Error while decoding hydra's response bytes"})
		}
	}
	panic(gohtypes.Error{Code: 500, Err: err, Message: "Error while communicating with Hydra"})
}
