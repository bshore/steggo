package encoders

import (
	"bytes"
	"encoding/ascii85"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// ==============================
// Rot13 Encoder/Decoder funcs
// ==============================
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

// Rot13 Encode and Decode are the same thing
func Rot13(input string) string {
	return strings.Map(rot13, input)
}

// ==============================
// Base16 Encoder/Decoder funcs
// ==============================

// Encode16 takes a message and hex encodes it
func Encode16(msg string) string {
	b := []byte(msg)
	encoded := make([]byte, len(b))
	_ = hex.Encode(encoded, b)
	return string(encoded)
}

// Decode16 takes an encoded hex string and decodes it
func Decode16(src string) (string, error) {
	b := []byte(src)
	decoded := make([]byte, len(b))
	_, err := hex.Decode(decoded, b)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// ==============================
// Base32 Encoder/Decoder funcs
// ==============================

// Encode32 takes in a message string and applies Base32 encoding
func Encode32(msg string) string {
	return base32.StdEncoding.EncodeToString([]byte(msg))
}

// Decode32 takes an encoded Base32 string and decodes it
func Decode32(src string) (string, error) {
	decoded, err := base32.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// ==============================
// Base64 Encoder/Decoder funcs
// ==============================

// Encode64 takes in a message string and Base64 encodes it
func Encode64(msg string) string {
	return base64.StdEncoding.EncodeToString([]byte(msg))
}

// Decode64 takes an encoded Base64 string and decodes it
func Decode64(src string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// ==============================
// Base85 Encoder/Decoder funcs
// ==============================

// Encode85 takes in a message string and Base85 encodes it
func Encode85(msg string) string {
	b := []byte(msg)
	encoded := make([]byte, ascii85.MaxEncodedLen(len(b)))
	_ = ascii85.Encode(encoded, b)
	return string(encoded)
}

// Decode85 takes an encoded Base85 string and decodes it
func Decode85(src string) (string, error) {
	b := []byte(src)
	decoded := make([]byte, len(b))
	numBytes, _, err := ascii85.Decode(decoded, b, true)
	if err != nil {
		return "", err
	}
	decoded = decoded[:numBytes]
	// remove /x00 null bytes before returning
	decoded = bytes.Trim(decoded, "\x00")
	return string(decoded), nil
}
