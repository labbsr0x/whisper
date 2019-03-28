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

// AcceptLoginRequestPayload holds the data to communicate with hydra's accept login api
type AcceptLoginRequestPayload struct {
	Subject     string `json:"subject"`
	Remember    bool   `json:"remember"`
	RememberFor int    `json:"remember_for"`
	ACR         string `json:"acr"`
}

// AcceptConsentRequestPayload holds the data to communicate with hydra's accept consent api
type AcceptConsentRequestPayload struct {
	GrantScope               []string            `json:"grant_scope"`
	GrantAccessTokenAudience interface{}         `json:"requested_access_token_audience"`
	Remember                 bool                `json:"remember"`
	RememberFor              int                 `json:"remember_for"`
	Session                  TokenSessionPayload `json:"session"`
}

// TokenSessionPayload holds additional data to be carried with the created token
type TokenSessionPayload struct {
	Owner string `json:"owner"`
}

// RejectConsentRequestPayload holds the data to communicate with hydra's reject consent api
type RejectConsentRequestPayload struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
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
func (hydra *HydraClient) AcceptLoginRequest(challenge string, payload AcceptLoginRequestPayload) map[string]interface{} {
	data, _ := json.Marshal(&payload)
	return hydra.put("login", "accept", challenge, data)
}

// GetConsentRequestInfo retrieves information to drive decisions over how to deal with the consent request
func (hydra *HydraClient) GetConsentRequestInfo(challenge string) map[string]interface{} {
	return hydra.get("consent", challenge)
}

// AcceptConsentRequest sends an accept login request to hydra
func (hydra *HydraClient) AcceptConsentRequest(challenge string, payload AcceptConsentRequestPayload) map[string]interface{} {
	data, _ := json.Marshal(&payload)
	return hydra.put("consent", "accept", challenge, data)
}

// RejectConsentRequest sends a reject login request to hydra
func (hydra *HydraClient) RejectConsentRequest(challenge string, payload RejectConsentRequestPayload) map[string]interface{} {
	data, _ := json.Marshal(&payload)
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
