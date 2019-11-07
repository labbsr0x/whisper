package misc

import (
	"fmt"
	"strings"
)

const (
	PasswordMinCharacters       = 12
	PasswordMaxCharacters       = 30
	PasswordMinUniqueCharacters = 7
)

func ValidatePassword(password, username, email string) error {
	if len(password) < PasswordMinCharacters {
		return fmt.Errorf("Your password should have at least %v characters", PasswordMinCharacters)
	}

	if len(password) > PasswordMaxCharacters {
		return fmt.Errorf("Your password should have at most %v characters", PasswordMaxCharacters)
	}

	pass := strings.ToLower(password)
	user := strings.ToLower(username)
	mail := strings.ToLower(email)

	if strings.Contains(pass, user) || strings.Contains(user, pass) {
		return fmt.Errorf("Your password is too similar to your username")
	}

	if strings.Contains(pass, mail) || strings.Contains(mail, pass) {
		return fmt.Errorf("Your password is too similar to your email")
	}

	if CountUniqueCharacters(pass) < PasswordMinUniqueCharacters {
		return fmt.Errorf("Your password should have at least %v unique characters", PasswordMinUniqueCharacters)
	}

	return nil
}

func GetPasswordTooltip() string {
	return fmt.Sprintf("<div style=\"text-align: left;\"> Password Rules:<br> 1. At least %v characters<br> 2. At most %v characters<br> 3. At least %v unique characters<br> 4. Differ from username and email<br> </div>", PasswordMinCharacters, PasswordMaxCharacters, PasswordMinUniqueCharacters)
}
