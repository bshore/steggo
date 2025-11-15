package encoders

import (
	"bytes"
	"compress/gzip"
	"encoding/ascii85"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
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
	encoded := make([]byte, hex.EncodedLen(len(b)))
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

// ==============================
// Gzip Compress/Decompress funcs
// ==============================

// Gzip takes in a message string and Gzip compresses it, and then Base64 encodes it
// for embedding
func Gzip(msg []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer: %v", err)
	}
	_, err = writer.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %v", err)
	}
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %v", err)
	}
	out := Encode64(buf.String())
	return []byte(out), nil
}

// Gunzip takes in a Gzip compressed byte slice, Base64 decodes it and decompresses it
func Gunzip(data []byte) ([]byte, error) {
	bs, err := Decode64(string(data))
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(strings.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer reader.Close()
	return io.ReadAll(reader)
}
