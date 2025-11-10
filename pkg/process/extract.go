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

// ExtractMsgFromGif extracts data from a GIF by reading unused palette slots
// This mirrors the embedding process which only modifies unused colors
func ExtractMsgFromGif(file *gif.GIF) (*Header, []byte, error) {
	var headBytes, msgBytes []byte
	var headerFound bool
	var header = &Header{}

	// Analyze the GIF to find unused palette slots (same as during embedding)
	capacities := AnalyzeGifCapacity(file)

	// totalCapacity := CalculateTotalGifCapacity(capacities)
	// fmt.Printf("Extracting from GIF with %d data values capacity across %d frames\n", totalCapacity, len(file.Image))

	// Extract data from unused palette slots in each frame
	for _, capacity := range capacities {
		if len(capacity.UnusedIndices) == 0 {
			continue // No unused slots in this frame
		}

		frame := file.Image[capacity.FrameIndex]

		// Extract from each unused palette color
		for _, paletteIdx := range capacity.UnusedIndices {
			// Get the palette color
			r, g, b, _ := frame.Palette[paletteIdx].RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// Extract the byte using the same 2-3-3 bit pattern as PNG/BMP
			extractedByte := extractFromColor(r8, g8, b8)

			if headerFound {
				if len(msgBytes) < header.Size {
					msgBytes = append(msgBytes, extractedByte)
				} else {
					// All message bytes extracted
					return header, msgBytes, nil
				}
			} else {
				// Still looking for header
				headBytes = append(headBytes, extractedByte)
				if header.Found(headBytes) {
					headerFound = true
					// fmt.Printf("Header found: %d bytes, type: %s, encoding: %s\n", header.Size, header.SrcType, header.PreEncoding)
				}
			}
		}
	}

	if !headerFound {
		return nil, nil, fmt.Errorf("failed to extract message: header not found")
	}

	if len(msgBytes) < header.Size {
		return nil, nil, fmt.Errorf("failed to extract complete message: got %d bytes, expected %d", len(msgBytes), header.Size)
	}

	return header, msgBytes, nil
}
