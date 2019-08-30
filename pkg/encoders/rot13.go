package encoders

import "strings"

func rot13(x rune) rune {
	capital := x >= 'A' && x <= 'Z'
	if !capital && (x < 'a' || x > 'z') {
		return x // Not a letter
	}
	x += 13
	if capital && x > 'Z' || !capital && x > 'z' {
		x -= 26
	}
	return x
}

// Rot13 is called inside strings.Map(Rot13, input)
// Rot13 encode & decode are the same thing
func Rot13(input string) string {
	return strings.Map(rot13, input)
}
