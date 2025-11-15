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

func ExtractMsgFromGIF(file *gif.GIF) (*Header, []byte, error) {
	var headBytes, msgBytes []byte
	var headerFound bool
	var header = &Header{}

	for frameIdx := range file.Image {
		for paletteIdx := 0; paletteIdx < len(file.Image[frameIdx].Palette); paletteIdx++ {
			r, g, b, _ := file.Image[frameIdx].Palette[paletteIdx].RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			extractedByte := extractFromColor(r8, g8, b8)
			if extractedByte == 0x00 {
				// Always skip zero bytes.
				// Embedder does not embed in zero bytes
				// because it leads to message corruption.
				continue
			}

			if headerFound {
				if len(msgBytes) < header.Size {
					msgBytes = append(msgBytes, extractedByte)
				} else {
					return header, msgBytes, nil
				}
			} else {
				headBytes = append(headBytes, extractedByte)
				if header.Found(headBytes) {
					headerFound = true
				}
			}
		}
	}

	if !headerFound {
		return nil, nil, fmt.Errorf("failed to extract message, header not found: %s", string(headBytes[:20]))
	}

	if len(msgBytes) < header.Size {
		return nil, nil, fmt.Errorf("failed to extract complete message: got %d bytes, expected %d", len(msgBytes), header.Size)
	}

	return header, msgBytes, nil
}
