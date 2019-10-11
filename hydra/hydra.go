package hydra

import (
	"encoding/json"
	"github.com/labbsr0x/goh/gohclient"
	"github.com/labbsr0x/goh/gohtypes"
	"net/http"
)

type Api interface {
	GetLoginRequestInfo(challenge string) map[string]interface{}
	AcceptLoginRequest(challenge string, payload AcceptLoginRequestPayload) map[string]interface{}
	GetConsentRequestInfo(challenge string) map[string]interface{}
	AcceptConsentRequest(challenge string, payload AcceptConsentRequestPayload) map[string]interface{}
	RejectConsentRequest(challenge string, payload RejectConsentRequestPayload) map[string]interface{}
}

type DefaultHydraHelper struct {
	client *gohclient.Default
}

func (dhh *DefaultHydraHelper) Init (adminURL string) Api {
	var err error

	dhh.client, err = gohclient.New(nil, adminURL)
	gohtypes.PanicIfError("Unable to create the client", http.StatusBadRequest, err)

	return dhh
}

// GetLoginRequestInfo retrieves information to drive decisions over how to deal with the login request
func (dhh *DefaultHydraHelper) GetLoginRequestInfo(challenge string) map[string]interface{} {
	return get(dhh.client, "login", challenge)
}

// AcceptLoginRequest sends an accept login request to hydra
func (dhh *DefaultHydraHelper) AcceptLoginRequest(challenge string, payload AcceptLoginRequestPayload) map[string]interface{} {
	data, _ := json.Marshal(&payload)
	return put(dhh.client, "login", challenge, "accept", data)
}

// GetConsentRequestInfo retrieves information to drive decisions over how to deal with the consent request
func (dhh *DefaultHydraHelper) GetConsentRequestInfo(challenge string) map[string]interface{} {
	return get(dhh.client, "consent", challenge)
}

// AcceptConsentRequest sends an accept login request to hydra
func (dhh *DefaultHydraHelper) AcceptConsentRequest(challenge string, payload AcceptConsentRequestPayload) map[string]interface{} {
	data, _ := json.Marshal(&payload)
	return put(dhh.client, "consent", challenge, "accept", data)
}

// RejectConsentRequest sends a reject login request to hydra
func (dhh *DefaultHydraHelper) RejectConsentRequest(challenge string, payload RejectConsentRequestPayload) map[string]interface{} {
	data, _ := json.Marshal(&payload)
	return put(dhh.client, "consent", challenge, "reject", data)
}