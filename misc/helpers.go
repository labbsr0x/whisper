package misc

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
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

// GenerateSalt a salt string with 16 bytes of crypto/rand data.
func GenerateSalt() string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(randomBytes)
}

// GetEncryptedPassword builds an encrypted password with hmac(sha512)
func GetEncryptedPassword(secretKey, password, salt string) string {
	hash := hmac.New(sha512.New, []byte(secretKey))
	io.WriteString(hash, password+salt)
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
