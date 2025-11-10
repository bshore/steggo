package process

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
)

// EmbedMsgInImage takes the message string and embeds it
// in the source file's byte string using Least Significant Bit(s)
func EmbedMsgInImage(data []byte, file image.Image) (draw.Image, error) {
	var bitsIndex int
	var newR, newG, newB uint16
	bounds := file.Bounds()
	pixels := bounds.Max.X * bounds.Max.Y
	if !(len(data) < pixels) {
		return nil, fmt.Errorf("message won't fit in image: %v bits to embed, %v pixels available", len(data), pixels)
	}
	newFile := image.NewNRGBA64(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := file.At(x, y).RGBA()
			// If the iteration is still under the length of message bits
			if bitsIndex < len(data) {
				newR = embedIn16BitColor(data[bitsIndex], r)
				// Check if there is a next bit pair to embed
				if bitsIndex+1 < len(data) {
					newG = embedIn16BitColor(data[bitsIndex+1], g)
					// Check if there is a next bit pair to embed
					if bitsIndex+2 < len(data) {
						newB = embedIn16BitColor(data[bitsIndex+2], b)
					} else {
						// No more message bits to embed, copy color value
						newB = uint16(b)
					}
				} else {
					// No more message bits to embed, copy color value
					newG = uint16(g)
				}
			} else {
				// No more message bits to embed, just copy the rest of the pixels
				newR = uint16(r)
				newG = uint16(g)
				newB = uint16(b)
			}
			newColor := color.NRGBA64{
				R: newR,
				G: newG,
				B: newB,
				A: uint16(a),
			}
			newFile.SetNRGBA64(x, y, newColor)
			bitsIndex += 3
		}
	}
	return newFile, nil
}

// EmbedMsgInGIF embeds data into a GIF using only unused palette slots
// This ensures zero visual artifacts since unused colors are never displayed
func EmbedMsgInGIF(data []byte, file *gif.GIF) (*gif.GIF, error) {
	// Analyze capacity
	capacities := AnalyzeGifCapacity(file)
	totalCapacity := CalculateTotalGifCapacity(capacities)

	// Check if message fits
	if len(data) > totalCapacity {
		return nil, fmt.Errorf("message won't fit: need %d data values, have %d capacity", len(data), totalCapacity)
	}

	// fmt.Printf("GIF capacity: %d data values across %d frames (approx %d message bytes)\n", totalCapacity, len(file.Image), totalCapacity/3)

	// Embed data frame by frame using unused palette slots
	dataIndex := 0
	for _, capacity := range capacities {
		if dataIndex >= len(data) {
			break // All data embedded
		}

		if len(capacity.UnusedIndices) == 0 {
			continue // No unused slots in this frame
		}

		frame := file.Image[capacity.FrameIndex]
		// fmt.Printf("Frame %d: %d unused palette slots\n", capacity.FrameIndex, len(capacity.UnusedIndices))

		// Embed data into unused palette colors
		// Note: data is already split into 2-3-3 format by FinalizeMessage
		// Each unused color can hold 3 data bytes (R, G, B components)
		for _, paletteIdx := range capacity.UnusedIndices {
			if dataIndex >= len(data) {
				break
			}

			// Get current color
			r, g, b, a := frame.Palette[paletteIdx].RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Embed R component (2 or 3 bits depending on 128 flag)
			if dataIndex < len(data) {
				r8 = embedInColor(data[dataIndex], r8)
				dataIndex++
			}

			// Embed G component (3 bits)
			if dataIndex < len(data) {
				g8 = embedInColor(data[dataIndex], g8)
				dataIndex++
			}

			// Embed B component (3 bits)
			if dataIndex < len(data) {
				b8 = embedInColor(data[dataIndex], b8)
				dataIndex++
			}

			// Update the palette color (safe because this index is unused)
			frame.Palette[paletteIdx] = color.RGBA{R: r8, G: g8, B: b8, A: a8}
		}
	}

	if dataIndex < len(data) {
		return nil, fmt.Errorf("failed to embed all data: embedded %d of %d bytes", dataIndex, len(data))
	}

	// fmt.Printf("Successfully embedded %d bytes\n", dataIndex)
	return file, nil
}
