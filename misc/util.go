package misc

func CountUniqueCharacters(str string) int {
	counter := make(map[int32]int)

	for _, value := range str {
		counter[value]++
	}

	return len(counter)
}
