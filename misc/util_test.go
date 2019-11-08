package misc

import "testing"

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
