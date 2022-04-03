package process

import (
	"fmt"
	"image"
	"image/gif"
)

// ExtractMsgFromImage takes an Image that has had a message embedded
// inside it and extracts the message using Least Significant Bit(s)
func ExtractMsgFromImage(file image.Image) (*Header, []byte, error) {
	var headBytes, msgBytes []byte
	var headerFound bool
	var header = &Header{}
	bounds := file.Bounds()
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := file.At(x, y).RGBA()
			if headerFound {
				if len(msgBytes) < header.Size {
					msgbyte := extractFromColor(uint8(r), uint8(g), uint8(b))
					msgBytes = append(msgBytes, msgbyte)
				}
				continue
			} else {
				// Build up headBytes until the full header is identified
				headbyte := extractFromColor(uint8(r), uint8(g), uint8(b))
				headBytes = append(headBytes, headbyte)
				if header.Found(headBytes) {
					headerFound = true
				}
			}
		}
	}

	if !headerFound {
		return nil, nil, fmt.Errorf("failed to extract message")
	}

	return header, msgBytes, nil
}

// ExtractMsgFromGif takes a GIF that has had a message embedded
// inside it and extracts the message using Least Significant Bit(s)
func ExtractMsgFromGif(file *gif.GIF) (*Header, []byte, error) {
	var headBytes, msgBytes []byte
	var headerFound bool
	var header = &Header{}
	// For each image frame
	for _, img := range file.Image {
		bounds := img.Bounds()
		// For each vertical row
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			// vor each pixel in each row
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				if headerFound {
					if len(msgBytes) < header.Size {
						msgbyte := extractFromColor(uint8(r), uint8(g), uint8(b))
						msgBytes = append(msgBytes, msgbyte)
					}
					continue
				} else {
					// Build up headBytes & headBits until it can be Unmarshaled
					headByte := extractFromColor(uint8(r), uint8(g), uint8(b))
					headBytes = append(headBytes, headByte)
					if header.Found(headBytes) {
						headerFound = true
					}
				}
			}
		}
	}

	if !headerFound {
		return nil, nil, fmt.Errorf("failed to extract message")
	}

	return header, msgBytes, nil
}
