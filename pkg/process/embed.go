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

// EmbedMsgInGIF takes the message data and embeds it into the GIF file's
// Local Color Palette.
func EmbedMsgInGIF(data []byte, file *gif.GIF) (*gif.GIF, error) {
	// Find all non-embedable colors in the the whole GIF
	var nonEmbedableColors int
	for frameIdx := range file.Image {
		for paletteIdx := 0; paletteIdx < len(file.Image[frameIdx].Palette); paletteIdx++ {
			r, g, b, _ := file.Image[frameIdx].Palette[paletteIdx].RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			extractedByte := extractFromColor(r8, g8, b8)
			if extractedByte == 0x00 {
				nonEmbedableColors++
			}
		}
	}
	// TODO:
	// Maybe need to add like a --gif-force option that will duplicate frames and scale the delay to fit the message?
	//
	//
	totalCapacity := (len(file.Image) * 256 * 3) - nonEmbedableColors
	if len(data) > totalCapacity {
		return nil, fmt.Errorf("message won't fit: need %d data values, have %d capacity", len(data), totalCapacity)
	}

	bitsIndex := 0
	for frameIdx := range file.Image {
		if bitsIndex >= len(data) {
			break
		}

		for paletteIdx := 0; paletteIdx < len(file.Image[frameIdx].Palette); paletteIdx++ {
			if bitsIndex >= len(data) {
				break
			}

			r, g, b, a := file.Image[frameIdx].Palette[paletteIdx].RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			extractedByte := extractFromColor(r8, g8, b8)
			if extractedByte == 0x00 {
				// Always skip bytes that produce a zero byte,
				// standard library GIF encoder optimizes by LSB
				// and will destroy any LSB modificaiton we make,
				// setting them back to 0x00
				continue
			}

			if bitsIndex < len(data) {
				r8 = embedInColor(data[bitsIndex], r8)
				if bitsIndex+1 < len(data) {
					g8 = embedInColor(data[bitsIndex+1], g8)
					if bitsIndex+2 < len(data) {
						b8 = embedInColor(data[bitsIndex+2], b8)
					}
				}
			}

			file.Image[frameIdx].Palette[paletteIdx] = color.RGBA{R: r8, G: g8, B: b8, A: a8}
			bitsIndex += 3
		}
	}

	if bitsIndex < len(data) {
		return nil, fmt.Errorf("failed to embed all data: embedded %d of %d bytes", bitsIndex, len(data))
	}

	return file, nil
}
