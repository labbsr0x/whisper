package misc

import (
	"fmt"
	"strings"
)

const (
	// PasswordMinChar is the minimum number of characters the password should have
	PasswordMinChar = 12

	// PasswordMaxChar is the maximum number of characters the password should have
	PasswordMaxChar = 30

	// PasswordMinUniqueChar is the minimum number of unique characters the password should have
	PasswordMinUniqueChar = 7
)

var passwordTooltipTemplate = `
<div style=\"text-align: left;\"> 
	Password Rules:<br> 
	1. At least %v characters<br> 
	2. At most %v characters<br> 
	3. At least %v unique characters<br> 
	4. Differ from username and email<br> 
</div>
`

// PasswordTooltip is an explanation snippet of the rules to a valid password
var PasswordTooltip = fmt.Sprintf(passwordTooltipTemplate, PasswordMinChar, PasswordMaxChar, PasswordMinUniqueChar)

// ValidatePassword verify if a password is valid
func ValidatePassword(password, username, email string) error {
	if len(password) < PasswordMinChar {
		return fmt.Errorf("password should have at least %v characters", PasswordMinChar)
	}

	if len(password) > PasswordMaxChar {
		return fmt.Errorf("password should have at most %v characters", PasswordMaxChar)
	}

	pass := strings.ToLower(password)
	user := strings.ToLower(username)
	mail := strings.ToLower(email)

	if strings.Contains(pass, user) || strings.Contains(user, pass) {
		return fmt.Errorf("password is too similar to your username")
	}

	if strings.Contains(pass, mail) || strings.Contains(mail, pass) {
		return fmt.Errorf("password is too similar to your email")
	}

	if CountUniqueCharacters(pass) < PasswordMinUniqueChar {
		return fmt.Errorf("password should have at least %v unique characters", PasswordMinUniqueChar)
	}

	return nil
}
