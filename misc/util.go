package misc

import (
	"fmt"
	"regexp"
)

func CountUniqueCharacters(str string) int {
	counter := make(map[int32]int)

	for _, value := range str {
		counter[value]++
	}

	return len(counter)
}

func VerifyEmail(email string) error {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(email) {
		return fmt.Errorf("Invalid email")
	}

	return nil
}
