package misc

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type IPayload interface {
	Check() error
}

// UnmarshalPayloadFromRequest initializes payload from an http request form and triggers the check function of the payload
func UnmarshalPayloadFromRequest(p IPayload, r *http.Request) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("Not possible to parse update payload")
	}

	err = json.Unmarshal(data, &p)
	if err != nil {
		return fmt.Errorf("Not possible to unmarshal update payload")
	}

	logrus.Debugf("Payload: '%v'", p)

	return p.Check()
}
