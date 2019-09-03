package process

import (
	"strconv"
	"strings"
)

// ZeroPadLeft left pads a string with zeros until the string is
// 8 characters long
func ZeroPadLeft(str string) string {
	if len(str) == 8 {
		return str
	}
	for {
		str = "0" + str
		if len(str) == 8 {
			return str
		}
	}
}

// BreakupMessageBytes busts apart each message byte into
// string arrays of two for swapping Least Significant Bits
// of the image file's writable bytes
func BreakupMessageBytes(msg string) []string {
	msgBytes := []byte(msg)
	var bitArr []string
	for _, b := range msgBytes {
		// Kinda dumb, I'm sure there's a better way to
		// get a binary representation of a byte...
		binStr := strconv.FormatInt(int64(b), 2)
		bits := strings.Split(ZeroPadLeft(binStr), "")
		for _, bit := range bits {
			bitArr = append(bitArr, bit)
		}
	}
	return bitArr
}
