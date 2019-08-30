package encoders

import "fmt"

// EncType enum for handling encoding types
type EncType int

const (
	// R13 short for Rot13 EncType
	R13 EncType = iota
	// B16 short for Base16 EncType
	B16 EncType = iota
	// B32 short for Base16 EncType
	B32 EncType = iota
	// B64 short for Base16 EncType
	B64 EncType = iota
	// B85 short for Base16 EncType
	B85 EncType = iota
)

var encMap = map[int]string{
	0: "rot13",
	1: "b16",
	2: "b32",
	3: "b64",
	4: "b85",
}

func (e EncType) String() string {
	return encMap[int(e)]
}

// EncTypeFromString takes a string and returns an EncType
func EncTypeFromString(s string) (EncType, error) {
	for e, str := range encMap {
		if s == str {
			return EncType(e), nil
		}
	}
	return 0, fmt.Errorf("unknown encoding type: %v", s)
}

// ApplyPreEncoding encodes the message with each type of encoding passed through cli
func ApplyPreEncoding(msg string, encs []EncType) string {
	// Apply the pre encoding junk
	return ""
}
