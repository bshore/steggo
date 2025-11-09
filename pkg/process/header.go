package process

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bshore/steggo/pkg/encoders"
)

// Header is a prefix to to identify information used during extraction.
type Header struct {
	Size        int
	SrcType     string
	PreEncoding string
}

// Found checks if bytes has the header string termination characters !/
func (h *Header) Found(b []byte) bool {
	if strings.HasSuffix(string(b), "!/") {
		headerPieces := strings.Split(strings.TrimSuffix(string(b), "!/"), ",")
		size, _ := strconv.ParseInt(headerPieces[0], 10, 64)
		h.Size = int(size)
		h.SrcType = headerPieces[1]
		h.PreEncoding = headerPieces[2]
		return true
	}
	return false
}

// NewBytesHeader takes the source parameters and returns a minimal representation to identify
//   the embedded data so we can later extract it back out
//
// Example:
// - "1024,0,1/2!/"
//   - 1024 indicates that the embedded message has a length of 1024 characters
//   - 0 indicates that the source type was a png
//   - 1/2!/ indicates that the message was pre-encoded with b16 and b32 before embed
func NewHeaderBytes(input, srcType string, encoders []encoders.EncType) []byte {
	header := Header{
		Size:    len(input),
		SrcType: srcType,
	}

	if len(encoders) == 0 {
		header.PreEncoding = "!/"
	} else {
		var encStrs []string
		for i := range encoders {
			encStrs = append(encStrs, encoders[i].String())
		}
		header.PreEncoding = fmt.Sprintf("%s!/", strings.Join(encStrs, "/"))
	}
	return []byte(fmt.Sprintf("%d,%s,%s", header.Size, header.SrcType, header.PreEncoding))
}

// FinalizeMessage transforms the header and message into it's final form for R,G,B least significant bit insertion
func FinalizeMessage(header []byte, msg string) []byte {
	msgBytes := []byte(string(header) + msg)
	var bitArr []byte
	for _, b := range msgBytes {
		// Get bit values in a group of 2-3-3 (R-G-B)
		// sevenEight uses & 131 to set the 128 bit, so embedding knows to zero out
		// the last 2 bits of a color value, instead of zeroing out the last 3 bits
		sevenEight := (b >> 6) & 131 // shifts bb------ to ------bb and gets last 2 bits value
		fourFiveSix := (b >> 3) & 7  // shifts --bbb--- to -----bbb and gets last 3 bits value
		oneTwoThree := b & 7         // just gets -----bbb last 3 bits value
		bitArr = append(bitArr, uint8(sevenEight), uint8(fourFiveSix), uint8(oneTwoThree))
	}
	return bitArr
}
