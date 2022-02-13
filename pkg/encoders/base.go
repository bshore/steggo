package encoders

import (
	"fmt"
	"strings"
)

// EncType enum for handling encoding types
type EncType int

const (
	// R13 short for Rot13 EncType
	R13 EncType = iota
	// B16 short for Base16 EncType
	B16 EncType = iota
	// B32 short for Base32 EncType
	B32 EncType = iota
	// B64 short for Base64 EncType
	B64 EncType = iota
	// B85 short for Base85 EncType
	B85 EncType = iota
)

var encMap = map[int]string{
	0: "r13",
	1: "b16",
	2: "b32",
	3: "b64",
	4: "b85",
}

func (e EncType) String() string {
	return fmt.Sprintf("%d", e)
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

func FromStrSlice(encStrSlice []string) ([]EncType, string) {
	var out []EncType
	var errs []string
	for i := range encStrSlice {
		enc, err := EncTypeFromString(encStrSlice[i])
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			out = append(out, enc)
		}
	}
	if len(errs) == 0 {
		return out, ""
	}
	return out, strings.Join(errs, "\n")
}

// ApplyPreEncoding encodes the message with each type of encoding passed through cli
func ApplyPreEncoding(msg string, encs []EncType) string {
	for _, enc := range encs {
		if enc == R13 {
			msg = Rot13(msg)
		} else if enc == B16 {
			msg = Encode16(msg)
		} else if enc == B32 {
			msg = Encode32(msg)
		} else if enc == B64 {
			msg = Encode64(msg)
		} else if enc == B85 {
			msg = Encode85(msg)
		}
	}
	return msg
}
