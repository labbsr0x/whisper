package misc

import (
	"fmt"
	"github.com/labbsr0x/goh/gohtypes"
	"net/http"
	"strings"
)

const (
	passMinCharacters       = 12
	passMaxCharacters       = 30
	passMinUniqueCharacters = 7
)

func ValidatePassword (password, username, email string) {
	if len(password) < passMinCharacters {
		gohtypes.Panic(fmt.Sprintf("Your password should have at least %v characters", passMinCharacters), http.StatusBadRequest)
	}

	if len(password) > passMaxCharacters {
		gohtypes.Panic(fmt.Sprintf("Your password should have at most %v characters", passMaxCharacters), http.StatusBadRequest)
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

	if CountUniqueCharacters(pass) < passMinUniqueCharacters {
		gohtypes.Panic(fmt.Sprintf("your password should have at least %v unique characters", passMinUniqueCharacters), http.StatusBadRequest)
	}
}

func GetPasswordTooltip () string {
	return fmt.Sprintf("<div style=\"text-align: left;\"> Password Rules:<br> 1. At least %v characters<br> 2. At most %v characters<br> 3. At least %v unique characters<br> 4. Differ from username and email<br> </div>", passMinCharacters, passMaxCharacters, passMinUniqueCharacters)
}