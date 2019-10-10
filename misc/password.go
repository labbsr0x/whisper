package misc

import (
	"fmt"
	"github.com/labbsr0x/goh/gohtypes"
	"net/http"
	"strings"
)

const (
	PasswordMinCharacters       = 12
	PasswordMaxCharacters       = 30
	PasswordMinUniqueCharacters = 7
)

func ValidatePassword(password, username, email string) {
	if len(password) < PasswordMinCharacters {
		gohtypes.Panic(fmt.Sprintf("Your password should have at least %v characters", PasswordMinCharacters), http.StatusBadRequest)
	}

	if len(password) > PasswordMaxCharacters {
		gohtypes.Panic(fmt.Sprintf("Your password should have at most %v characters", PasswordMaxCharacters), http.StatusBadRequest)
	}

	pass := strings.ToLower(password)
	user := strings.ToLower(username)
	mail := strings.ToLower(email)

	if strings.Contains(pass, user) || strings.Contains(user, pass) {
		gohtypes.Panic("Your password is too similar to your username", http.StatusBadRequest)
	}

	if strings.Contains(pass, mail) || strings.Contains(mail, pass) {
		gohtypes.Panic("Your password is too similar to your email", http.StatusBadRequest)
	}

	if CountUniqueCharacters(pass) < PasswordMinUniqueCharacters {
		gohtypes.Panic(fmt.Sprintf("your password should have at least %v unique characters", PasswordMinUniqueCharacters), http.StatusBadRequest)
	}
}

func GetPasswordTooltip() string {
	return fmt.Sprintf("<div style=\"text-align: left;\"> Password Rules:<br> 1. At least %v characters<br> 2. At most %v characters<br> 3. At least %v unique characters<br> 4. Differ from username and email<br> </div>", PasswordMinCharacters, PasswordMaxCharacters, PasswordMinUniqueCharacters)
}
