package process

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"strconv"
	"strings"
)

// EmbedMsgInImage takes the message string and embeds it
// in the source file's byte string using Least Significant Bit
func EmbedMsgInImage(msg, format string, file image.Image) (draw.Image, error) {
	var bitsIndex int
	var err error
	var newR, newG, newB uint8
	bitArr := BreakupMessageBytes(msg)
	bitMax := len(bitArr) - 1
	bounds := file.Bounds()
	newFile := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := file.At(x, y).RGBA()
			// If the iteration is still under the length of message bits
			if bitsIndex < bitMax {
				newR, err = embedInColor(bitArr[bitsIndex], uint8(r))
				if err != nil {
					return nil, err
				}
				// Check if there is a next bit to embed
				if bitsIndex+1 < bitMax {
					newG, err = embedInColor(bitArr[bitsIndex+1], uint8(g))
					if err != nil {
						return nil, err
					}
					// Check if there is a next bit to embed
					if bitsIndex+2 < bitMax {
						newB, err = embedInColor(bitArr[bitsIndex+2], uint8(b))
						if err != nil {
							return nil, err
						}
						// Check if there is a next bit to embed
						// if bitsIndex+3 < bitMax {
						// 	newA, err = embedInColor(bitArr[bitsIndex+3], uint8(a))
						// 	if err != nil {
						// 		return nil, err
						// 	}
						// } else {
						// 	// No more message bits to embed, copy color value
						// 	newA = uint8(a)
						// }
					} else {
						// No more message bits to embed, copy color value
						newB = uint8(b)
					}
				} else {
					// No more message bits to embed, copy color value
					newG = uint8(g)
				}
			} else {
				// No more message bits to embed, just copy the rest of the pixels
				newR = uint8(r)
				newG = uint8(g)
				newB = uint8(b)
				// newA = uint8(a)
			}
			newColor := color.RGBA{
				R: newR,
				G: newG,
				B: newB,
				A: uint8(a),
			}
			newFile.SetRGBA(x, y, newColor)
			bitsIndex = bitsIndex + 4
		}
	}
	return newFile, nil
}

// EmbedMsgInGIF takes the message string and embeds it into a GIF file frame by frame
func EmbedMsgInGIF(msg, format string, file *gif.GIF) (draw.Image, error) {
	// Oof, this is gonna be fun
	return nil, nil
}

// There's gotta be a better way...
func embedInColor(a string, b uint8) (uint8, error) {
	colorVal := strconv.FormatUint(uint64(b), 2)
	colorBits := strings.Split(ZeroPadLeft(colorVal), "")
	// Last Bit becomes the bit of encoded msg
	colorBits[len(colorBits)-1] = a

	// Rejoin & reparse the new color values
	colorStr := strings.Join(colorBits, "")
	newColor, err := strconv.ParseUint(colorStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint8(newColor), nil
}
