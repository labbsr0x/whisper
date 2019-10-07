package misc

import (
	"fmt"
	"strings"
)

func ValidatePassword (password, username, email string) error {
	const (
		minCharacters = 12
		maxCharacters = 30
		minUnique = 7
	)

	if len(password) < minCharacters {
		return fmt.Errorf("your password should have at least %v characters", minCharacters)
	}

	if len(password) > maxCharacters {
		return fmt.Errorf("your password should have at most %v characters", maxCharacters)
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

	if CountUniqueCharacters(pass) < minUnique {
		return fmt.Errorf("your password should have at least %v unique characters", minUnique)
	}

	return nil
}
