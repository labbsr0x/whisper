package misc

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/abilioesteves/goh/gohtypes"
	"github.com/sirupsen/logrus"

	"github.com/abilioesteves/goh/gohclient"
)

// HydraClient holds data and methods to communicate with an hydra service instance
type HydraClient struct {
	BaseURL    *url.URL
	HTTPClient *gohclient.Default
}

// Init initializes a hydra client
func (hydra *HydraClient) Init(hydraEndpoint string) *HydraClient {
	hydra.BaseURL = new(url.URL)
	hydra.HTTPClient = gohclient.New("application/json")

	hydra.BaseURL.Host = hydraEndpoint
	hydra.BaseURL.Path = "/oauth2/auth/requests/"

	return hydra
}

// GetLoginRequestInfo retrieves information to drive decisions over how to deal with the login request
func (hydra *HydraClient) GetLoginRequestInfo(challenge string) map[string]interface{} {
	return hydra.get("login", challenge)
}

// AcceptLoginRequest sends an accept login request to hydra
func (hydra *HydraClient) AcceptLoginRequest(challenge string, data []byte) map[string]interface{} {
	return hydra.put("login", "accept", challenge, data)
}

// RejectLoginRequest sends a reject login request to hydra
func (hydra *HydraClient) RejectLoginRequest(challenge string, data []byte) map[string]interface{} {
	return hydra.put("login", "reject", challenge, data)
}

// GetConsentRequestInfo retrieves information to drive decisions over how to deal with the consent request
func (hydra *HydraClient) GetConsentRequestInfo(challenge string) map[string]interface{} {
	return hydra.get("consent", challenge)
}

// AcceptConsentRequest sends an accept login request to hydra
func (hydra *HydraClient) AcceptConsentRequest(challenge string, data []byte) map[string]interface{} {
	return hydra.put("consent", "accept", challenge, data)
}

// RejectConsentRequest sends a reject login request to hydra
func (hydra *HydraClient) RejectConsentRequest(challenge string, data []byte) map[string]interface{} {
	return hydra.put("consent", "reject", challenge, data)
}

func (hydra *HydraClient) get(flow, challenge string) map[string]interface{} {
	return hydra.treatResponse(hydra.HTTPClient.Get(path.Join(hydra.BaseURL.String(), flow, challenge)))
}

func (hydra *HydraClient) put(flow, challenge, action string, data []byte) map[string]interface{} {
	return hydra.treatResponse(hydra.HTTPClient.Put(path.Join(hydra.BaseURL.String(), flow, challenge, action), data))
}

func (hydra *HydraClient) treatResponse(resp *http.Response, data []byte, err error) map[string]interface{} {
	status, _ := strconv.Atoi(resp.Status)
	if status < 200 || status > 302 || err != nil {
		panic(gohtypes.Error{Code: status, Err: err, Message: "Error while communicating with Hydra"})
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		panic(gohtypes.Error{Code: 500, Err: err, Message: "Error while decoding hydra's response bytes"})
	}

	logrus.Infof("Result: %v", result)

	return result
}
