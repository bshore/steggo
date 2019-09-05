package process

import (
	"strconv"
	"strings"
)

// ZeroPadLeft left pads a string with zeros until the string is
// 8 characters long
func ZeroPadLeft(str string, length int) string {
	if len(str) == length {
		return str
	}
	for {
		str = "0" + str
		if len(str) == length {
			return str
		}
	}
}

// BreakupMessageBytes busts apart each message byte into
// string arrays of two for swapping Least Significant Bits
// of the image file's writable bytes
func BreakupMessageBytes(msg string) [][]string {
	msgBytes := []byte(msg)
	var bitArr [][]string
	for _, b := range msgBytes {
		// Kinda dumb, I'm sure there's a better way to
		// get a binary representation of a byte...
		binStr := strconv.FormatInt(int64(b), 2)
		bits := strings.Split(ZeroPadLeft(binStr, 8), "")
		for i, bit := range bits {
			// Current iteration is odd ?
			// Grab bits in pairs (01, 23, 45, 67)
			if i%2 != 0 && i != 0 {
				// 00000010 two's bit
				two := bit
				// 00000001 one's bit
				one := bits[i-1]
				bitArr = append(bitArr, []string{two, one})
			}
		}
	}
	return bitArr
}
