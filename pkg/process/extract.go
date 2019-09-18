package process

import (
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"strconv"
)

// ExtractMsgFromImage takes an Image that has had a message embedded
// inside it and extracts the message using Least Significant Bit(s)
func ExtractMsgFromImage(secret *Secret, file image.Image) (*Secret, error) {
	var err error
	var size int64
	var headBytes, msgBytes []byte
	var headerFound bool
	var header Header
	bounds := file.Bounds()
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := file.At(x, y).RGBA()
			if headerFound {
				if int64(len(msgBytes)) < size {
					msgbyte := extractFromColor(uint8(r), uint8(g), uint8(b))
					msgBytes = append(msgBytes, msgbyte)
				}
				continue
			} else {
				// Build up headBytes until it can be Unmarshaled
				headbyte := extractFromColor(uint8(r), uint8(g), uint8(b))
				headBytes = append(headBytes, headbyte)
				err = json.Unmarshal(headBytes, &header)
				if err == nil {
					headerFound = true
					size, err = strconv.ParseInt(header.Size, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("Error parsing secret size from header: (%v)", err)
					}
				}
			}
		}
	}
	secret.DataHeader = header
	secret.Message = msgBytes
	return secret, nil
}

// ExtractMsgFromGif takes a GIF that has had a message embedded
// inside it and extracts the message using Least Significant Bit(s)
func ExtractMsgFromGif(secret *Secret, file *gif.GIF) (*Secret, error) {
	var err error
	var size int64
	var headBytes, msgBytes []byte
	var headerFound bool
	var header Header
	// For each image frame
	for _, img := range file.Image {
		bounds := img.Bounds()
		// For each vertical row
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			// vor each pixel in each row
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				if headerFound {
					if int64(len(msgBytes)) < size {
						msgbyte := extractFromColor(uint8(r), uint8(g), uint8(b))
						msgBytes = append(msgBytes, msgbyte)
					}
					continue
				} else {
					// Build up headBytes until it can be Unmarshaled
					headByte := extractFromColor(uint8(r), uint8(g), uint8(b))
					headBytes = append(headBytes, headByte)
					err = json.Unmarshal(headBytes, &header)
					if err == nil {
						headerFound = true
						size, err = strconv.ParseInt(header.Size, 10, 64)
						if err != nil {
							return nil, fmt.Errorf("Error parsing secret size from header: (%v)", err)
						}
					}
				}
			}
		}
	}
	secret.DataHeader = header
	secret.Message = msgBytes
	return secret, nil
}

func extractFromColor(r, g, b uint8) byte {
	// Get last bits of each color to reconstruct a message byte
	rBits := r & 3
	gBits := g & 7
	bBits := b & 7

	var newByte uint8
	// Assign color bits and shift left for each color pixel
	newByte = newByte | (rBits & 3) // ------bb
	newByte = newByte << 3          // ---bb---
	newByte = newByte | (gBits & 7) // -----bbb
	newByte = newByte << 3          // bbbbb---
	newByte = newByte | (bBits & 7) // bbbbbbbb
	return newByte
}

// extract16BitColor(r, g, b uint32) byte {
//
//}
