package misc

import (
	"fmt"
	"strings"
)

const (
	passMinCharacters       = 12
	passMaxCharacters       = 30
	passMinUniqueCharacters = 7
)

func ValidatePassword (password, username, email string) error {
	if len(password) < passMinCharacters {
		return fmt.Errorf("your password should have at least %v characters", passMinCharacters)
	}

	if len(password) > passMaxCharacters {
		return fmt.Errorf("your password should have at most %v characters", passMaxCharacters)
	}

	pass := strings.ToLower(password)
	user := strings.ToLower(username)
	mail := strings.ToLower(email)

	if strings.Contains(pass, user) || strings.Contains(user, pass) {
		return fmt.Errorf("your password is too similar to your username")
	}

	if strings.Contains(pass, mail) || strings.Contains(mail, pass) {
		return fmt.Errorf("your password is too similar to your email")
	}

	if CountUniqueCharacters(pass) < passMinUniqueCharacters {
		return fmt.Errorf("your password should have at least %v unique characters", passMinUniqueCharacters)
	}

	return nil
}

func GetPasswordTooltip () string {
	return fmt.Sprintf("<div style=\"text-align: left;\"> Password Rules:<br> 1. At least %v characters<br> 2. At most %v characters<br> 3. At least %v unique characters<br> 4. Differ from username and email<br> </div>", passMinCharacters, passMaxCharacters, passMinUniqueCharacters)
}