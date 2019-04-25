package misc

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

// GetJSONStr gets the full json string from the toEncode structure
func GetJSONStr(toEncode interface{}) string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(toEncode)
	return buf.String()
}

// GetAccessTokenFromRequest is a helper method to recover an Access Token from a http request
func GetAccessTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	authURLParam := r.URL.Query().Get("token")
	var t string

	if len(authHeader) == 0 && len(authURLParam) == 0 {
		return "", fmt.Errorf("No Authorization Header or URL Param found")
	}

	if len(authHeader) > 0 {
		data := strings.Split(authHeader, " ")

		if len(data) != 2 {
			return "", fmt.Errorf("Bad Authorization Header")
		}

		t = data[0]

		if len(t) == 0 || t != "Bearer" {
			return "", fmt.Errorf("No Bearer Token found")
		}

		t = data[1]

	} else {
		t = authURLParam
	}

	if len(t) == 0 {
		return "", fmt.Errorf("Bad Authorization Token")
	}

	return t, nil
}

// GenerateSalt a salt string with 16 bytes of crypto/rand data.
func GenerateSalt() string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(randomBytes)
}

// GetEncryptedPassword builds an encrypted password with hmac(sha256)
func GetEncryptedPassword(secretKey, password, salt string) string {
	hash := hmac.New(sha256.New, []byte(secretKey))
	io.WriteString(hash, password+salt)
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
