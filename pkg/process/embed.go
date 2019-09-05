package process

import (
	"fmt"
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
	var newR, newG, newB, newA uint16
	bitArr := BreakupMessageBytes(msg)
	bitMax := len(bitArr) - 1
	bounds := file.Bounds()
	newFile := image.NewRGBA64(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := file.At(x, y).RGBA()
			// If the iteration is still under the length of message bits
			if bitsIndex < bitMax {
				newR, err = embedIn16BitColor(bitArr[bitsIndex], r)
				if err != nil {
					return nil, err
				}
				// Check if there is a next bit to embed
				if bitsIndex+1 < bitMax {
					newG, err = embedIn16BitColor(bitArr[bitsIndex+1], g)
					if err != nil {
						return nil, err
					}
					// Check if there is a next bit to embed
					if bitsIndex+2 < bitMax {
						newB, err = embedIn16BitColor(bitArr[bitsIndex+2], b)
						if err != nil {
							return nil, err
						}
						// Check if there is a next bit to embed
						if bitsIndex+3 < bitMax {
							newA, err = embedIn16BitColor(bitArr[bitsIndex+3], a)
							if err != nil {
								return nil, err
							}
						} else {
							// No more message bits to embed, copy color value
							newA = uint16(a)
						}
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
				newA = uint16(a)
			}
			newColor := color.RGBA64{
				R: newR,
				G: newG,
				B: newB,
				A: newA,
			}
			newFile.SetRGBA64(x, y, newColor)
			bitsIndex = bitsIndex + 4
		}
	}
	return newFile, nil
}

// EmbedMsgInGIF takes the message string and embeds it into a GIF file frame by frame
func EmbedMsgInGIF(msg, format string, file *gif.GIF) (*gif.GIF, error) {
	var bitsIndex int
	var err error
	var doneEmbedding bool
	var newR, newG, newB, newA uint8
	bitArr := BreakupMessageBytes(msg)
	bitMax := len(bitArr) - 1
	fmt.Println("Message bits to embed: ", bitMax)
	newGif := &gif.GIF{
		Image:           []*image.Paletted{},
		Delay:           file.Delay,
		LoopCount:       file.LoopCount,
		Disposal:        file.Disposal,
		Config:          file.Config,
		BackgroundIndex: file.BackgroundIndex,
	}
	// For each image frame
	for i, img := range file.Image {
		colorPalette := getGIFColorPallete(img)
		fmt.Printf("On frame %v of %v\n", i, len(file.Image))
		bounds := img.Bounds()
		newFrame := image.NewPaletted(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), colorPalette)
		// For each vertical row
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			// For each pixel in each row
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				// If the iteration is still under the length of message bits
				if bitsIndex < bitMax {
					newR, err = embedIn8BitColor(bitArr[bitsIndex], uint8(r))
					if err != nil {
						return nil, err
					}
					// Check if there is a next bit to embed
					if bitsIndex+1 < bitMax {
						newG, err = embedIn8BitColor(bitArr[bitsIndex+1], uint8(g))
						if err != nil {
							return nil, err
						}
						// Check if there is a next bit to embed
						if bitsIndex+2 < bitMax {
							newB, err = embedIn8BitColor(bitArr[bitsIndex+2], uint8(b))
							if err != nil {
								return nil, err
							}
							// Check if there is a next bit to embed
							if bitsIndex+3 < bitMax {
								newA, err = embedIn8BitColor(bitArr[bitsIndex+3], uint8(a))
								if err != nil {
									return nil, err
								}
							} else {
								// No more message bits to embed, copy color value
								newA = uint8(a)
							}
						} else {
							// No more message bits to embed, copy color value
							newB = uint8(b)
						}
					} else {
						// No more message bits to embed, copy color value
						newG = uint8(g)
					}
					newColor := color.RGBA{
						R: newR,
						G: newG,
						B: newB,
						A: newA,
					}
					newFrame.Set(x, y, newColor)
					bitsIndex = bitsIndex + 4
				} else {
					// No more message bits to embed, just copy the remaining pixels for frame
					newR = uint8(r)
					newG = uint8(g)
					newB = uint8(b)
					newA = uint8(a)
					newColor := color.RGBA{
						R: newR,
						G: newG,
						B: newB,
						A: newA,
					}
					newFrame.Set(x, y, newColor)
					bitsIndex = bitsIndex + 4
					if !doneEmbedding {
						fmt.Printf("Done embedding Frame %v at row Y: %v  col X: %v\n", i, y, x)
					}
					doneEmbedding = true
				}
			} // End x
		} // End y
		if doneEmbedding {
			// Append the current frame since it may contain parts of the message
			newGif.Image = append(newGif.Image, newFrame)
			// Append the next frame and every frame after it and return
			newGif.Image = append(newGif.Image, file.Image[i+1:]...)
			return newGif, nil
		}
		newGif.Image = append(newGif.Image, newFrame)
	}
	return newGif, nil
}

func embedIn8BitColor(a []string, b uint8) (uint8, error) {
	colorVal := strconv.FormatUint(uint64(b), 2)
	colorBits := strings.Split(ZeroPadLeft(colorVal, 8), "")
	// Least Significant Bits become bit pair of encoded msg
	colorBits[len(colorBits)-2] = a[0]
	colorBits[len(colorBits)-1] = a[1]

	// Rejoin & reparse the new color values
	colorStr := strings.Join(colorBits, "")
	newColor, err := strconv.ParseUint(colorStr, 2, 64)
	if err != nil {
		return 0, err
	}
	return uint8(newColor), nil
}

func embedIn16BitColor(a []string, b uint32) (uint16, error) {
	colorVal := strconv.FormatUint(uint64(b), 2)
	colorBits := strings.Split(ZeroPadLeft(colorVal, 16), "")
	// Least Significant Bits become bit pair of encoded msg
	colorBits[len(colorBits)-2] = a[0]
	colorBits[len(colorBits)-1] = a[1]

	// Rejoin & reparse the new color values
	colorStr := strings.Join(colorBits, "")
	newColor, err := strconv.ParseUint(colorStr, 2, 64)
	if err != nil {
		return 0, err
	}

	return uint16(newColor), nil
}

func getGIFColorPallete(img *image.Paletted) color.Palette {
	var colors color.Palette
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			colors = append(colors, img.At(x, y))
		}
	}
	return colors
}
