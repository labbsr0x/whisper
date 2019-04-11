package misc

import (
	"fmt"
	"net/http"
	"strings"
)

// ConvertInterfaceArrayToStringArray a helper method to perform conversions between []interface{} and []string
func ConvertInterfaceArrayToStringArray(toConvert []interface{}) []string {
	toReturn := make([]string, 0)
	if l := len(toConvert); l > 0 {
		toReturn = make([]string, l)
		for i := 0; i < l; i++ {
			toReturn[i] = toConvert[i].(string)
		}
	}
	return toReturn
}

// GetAccessTokenFromRequest is a helper method to recover an Access Token from a http request
func GetAccessTokenFromRequest(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	var t string

	if len(auth) == 0 {
		return "", fmt.Errorf("No Authorization Header found")
	}

	data := strings.Split(auth, " ")

	if len(data) != 2 {
		return "", fmt.Errorf("Bad Authorization Header")
	}

	t = data[0]

	if len(t) == 0 || t != "Bearer" {
		return "", fmt.Errorf("No Bearer Token found")
	}

	t = data[1]

	if len(t) == 0 {
		return "", fmt.Errorf("Bad Authorization Header")
	}

	return t, nil
}
