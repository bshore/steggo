package encoders

import (
	"bytes"
	"encoding/ascii85"
)

// Encode85 takes a source string and encodes it to Base85
func Encode85(value string) string {
	b := []byte(value)
	encoded := make([]byte, ascii85.MaxEncodedLen(len(b)))
	_ = ascii85.Encode(encoded, b)
	return string(encoded)
}

// Decode85 takes a Base85 encoded string and returns the source string
func Decode85(value string) (string, error) {
	b := []byte(value)
	decoded := make([]byte, len(b))
	nDecodedBytes, _, err := ascii85.Decode(decoded, b, true)
	if err != nil {
		return "", err
	}
	decoded = decoded[:nDecodedBytes]
	//ascii85 adds /x00 null bytes at the end
	decoded = bytes.Trim(decoded, "\x00")
	return string(decoded), nil
}
