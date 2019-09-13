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
func EmbedMsgInImage(secret *Secret, file image.Image) (draw.Image, error) {
	var bitsIndex int
	var newR, newG, newB, newA uint16
	bounds := file.Bounds()
	pixels := bounds.Max.X * bounds.Max.Y
	if !(secret.Size < pixels) {
		return nil, fmt.Errorf("Secret message won't fit in image: %v LSB's to embed, %v pixels available", secret.Size, pixels)
	}
	newFile := image.NewRGBA64(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := file.At(x, y).RGBA()
			// If the iteration is still under the length of message bits
			if bitsIndex < secret.Size {
				newR = embedIn16BitColor(secret.Data[bitsIndex], r)
				// Check if there is a next bit pair to embed
				if bitsIndex+1 < secret.Size {
					newG = embedIn16BitColor(secret.Data[bitsIndex+1], g)
					// Check if there is a next bit pair to embed
					if bitsIndex+2 < secret.Size {
						newB = embedIn16BitColor(secret.Data[bitsIndex+2], b)
						// Check if there is a next bit pair to embed
						if bitsIndex+3 < secret.Size {
							newA = embedIn16BitColor(secret.Data[bitsIndex+3], a)
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

// EmbedMsgInGIF takes the message string and embeds it into a GIF file
// frame by frame using Least Significant Bit(s)
func EmbedMsgInGIF(secret *Secret, file *gif.GIF) (*gif.GIF, error) {
	var bitsIndex int
	var doneEmbedding bool
	var newR, newG, newB, newA uint8
	bounds := file.Image[0].Bounds()
	// Get Bounds of first frame. Since GIF's cannot change size/resolution
	// this is a good way to estimate how many pixels we have available
	// for embedding.
	pixels := (bounds.Max.X * bounds.Max.Y) * len(file.Image)
	if !(secret.Size < pixels) {
		return nil, fmt.Errorf("Secret message won't fit in image: %v LSB's to embed, %v pixels available", secret.Size, pixels)
	}
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
		bounds := img.Bounds()
		newFrame := image.NewPaletted(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), img.Palette)
		// For each vertical row
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			// For each pixel in each row
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				// If the iteration is still under the length of message bits
				if bitsIndex < secret.Size {
					newR = embedIn8BitColor(secret.Data[bitsIndex], uint8(r))
					// Check if there is a next bit pair to embed
					if bitsIndex+1 < secret.Size {
						newG = embedIn8BitColor(secret.Data[bitsIndex+1], uint8(g))
						// Check if there is a next bit pair to embed
						if bitsIndex+2 < secret.Size {
							newB = embedIn8BitColor(secret.Data[bitsIndex+2], uint8(b))
							// Check if there is a next bit pair to embed
							if bitsIndex+3 < secret.Size {
								newA = embedIn8BitColor(secret.Data[bitsIndex+3], uint8(a))
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

func embedIn8BitColor(a byte, b uint8) uint8 {
	c := b | a
	return uint8(c)
	// colorVal := strconv.FormatUint(uint64(b), 2)
	// colorBits := strings.Split(ZeroPadLeft(colorVal, 8), "")
	// // Least Significant Bits become bit pair of encoded msg
	// colorBits[len(colorBits)-2] = a[0]
	// colorBits[len(colorBits)-1] = a[1]

	// // Rejoin & reparse the new color values
	// colorStr := strings.Join(colorBits, "")
	// newColor, err := strconv.ParseUint(colorStr, 2, 64)
	// if err != nil {
	// 	return 0, err
	// }
	// return uint8(newColor), nil
}

func embedIn16BitColor(a uint8, b uint32) uint16 {
	c := b | uint32(a)
	return uint16(c)
	// colorVal := strconv.FormatUint(uint64(b), 2)
	// colorBits := strings.Split(ZeroPadLeft(colorVal, 16), "")
	// // Least Significant Bits become bit pair of encoded msg
	// colorBits[len(colorBits)-2] = a[0]
	// colorBits[len(colorBits)-1] = a[1]

	// // Rejoin & reparse the new color values
	// colorStr := strings.Join(colorBits, "")
	// newColor, err := strconv.ParseUint(colorStr, 2, 64)
	// if err != nil {
	// 	return 0, err
	// }

	// return uint16(newColor), nil
}
