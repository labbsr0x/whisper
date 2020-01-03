package misc

import (
	"strings"
	"testing"
)

const (
	mockUsername = "username"
	mockEmail    = "anotheruserwithcreativity@mail.com"
)

var invalidPasswords = []string{
	"123456789", // too small
	"12345678910111213141516171819202122232425",   // too big
	mockUsername + "passwordtop123",               // similar to username
	mockEmail[:strings.IndexByte(mockEmail, '@')], // similar to email
	"11111111111111111",                           // not unique enough
}

func TestValidatePassword(t *testing.T) {
	for _, password := range invalidPasswords {
		err := ValidatePassword(password, mockUsername, mockEmail)
		if err == nil {
			t.Fail()
		}
	}
}

var testCountUniqueCharactersData = []struct {
	input  string
	output int
}{
	{"123456789", 9},
	{"123456799", 8},
	{"123556789", 8},
	{"127345689", 9},
	{"111111111", 1},
	{"111111112", 2},
	{"abcdefghi", 9},
	{"abcabcacb", 3},
}

func TestCountUniqueCharacters(t *testing.T) {
	for _, test := range testCountUniqueCharactersData {
		if CountUniqueCharacters(test.input) != test.output {
			t.Fail()
		}
	}
}
