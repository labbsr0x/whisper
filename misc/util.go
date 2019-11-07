package misc

import (
	"fmt"
	"regexp"
)

// CountUniqueCharacters counts the unique characters in a string
func CountUniqueCharacters(str string) int {
	counter := make(map[int32]int)

	for _, value := range str {
		counter[value]++
	}

	return len(counter)
}

// VerifyEmail verify if the string is actually an email
func VerifyEmail(email string) error {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(email) {
		return fmt.Errorf("invalid email")
	}

	return nil
}
