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
	var newR, newG, newB uint8
	bounds := file.Bounds()
	pixels := bounds.Max.X * bounds.Max.Y
	if !(secret.Size < pixels) {
		return nil, fmt.Errorf("Secret message won't fit in image: %v LSB's to embed, %v pixels available", secret.Size, pixels)
	}
	newFile := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := file.At(x, y).RGBA()
			// If the iteration is still under the length of message bits
			if bitsIndex < secret.Size {
				newR = embedInColor(secret.Data[bitsIndex], uint8(r))
				// Check if there is a next bit pair to embed
				if bitsIndex+1 < secret.Size {
					newG = embedInColor(secret.Data[bitsIndex+1], uint8(g))
					// Check if there is a next bit pair to embed
					if bitsIndex+2 < secret.Size {
						newB = embedInColor(secret.Data[bitsIndex+2], uint8(b))
					} else {
						// No more message bits to embed, copy color value
						newB = uint8(b)
					}
				} else {
					// No more message bits to embed, copy color value
					newG = uint8(g)
				}
				// fmt.Printf("R: %08b\tG: %08b\tB: %08b\tA: %08b\n", uint8(newR), uint8(newG), uint8(newB), uint8(a))
			} else {
				// No more message bits to embed, just copy the rest of the pixels
				newR = uint8(r)
				newG = uint8(g)
				newB = uint8(b)
			}
			newColor := color.RGBA{
				R: newR,
				G: newG,
				B: newB,
				A: uint8(a),
			}
			newFile.SetRGBA(x, y, newColor)
			bitsIndex += 3
		}
	}
	return newFile, nil
}

// EmbedMsgInGIF takes the message string and embeds it into a GIF file
// frame by frame using Least Significant Bit(s)
func EmbedMsgInGIF(secret *Secret, file *gif.GIF) (*gif.GIF, error) {
	var previous, bitsIndex int
	var doneEmbedding bool
	var newR, newG, newB uint8
	var newFrame *image.Paletted
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
		// The image rectangle bounds
		bounds := img.Bounds()
		// An empty frame with the same size as the source GIF and an empty color palette
		newFrame = image.NewPaletted(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), []color.Color{})
		// For each vertical row
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			// For each pixel in each row
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				// If bitsIndex is divisible by 768, we can no longer embed in this frame
				if bitsIndex != 0 && bitsIndex%768 == 0 {
					newFrame.Set(x, y, img.At(x, y))
					continue
				} else {
					r, g, b, a := img.At(x, y).RGBA()
					// If the iteration is still under the length of message bits
					if bitsIndex < secret.Size {
						newR = embedInColor(secret.Data[bitsIndex], uint8(r))
						// Check if next msg byte to embed and if byte will fit
						if bitsIndex+1 < secret.Size && bitsIndex+1%768 != 0 {
							newG = embedInColor(secret.Data[bitsIndex+1], uint8(g))
							// Check if next msg byte to embed and if byte will fit
							if bitsIndex+2 < secret.Size && bitsIndex+1%768 != 0 {
								newB = embedInColor(secret.Data[bitsIndex+2], uint8(b))
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
							A: uint8(a),
						}
						newFrame.Set(x, y, newColor)
						bitsIndex += 3
					} else {
						// No more message bits to embed, just copy the remaining pixels for frame
						newColor := color.RGBA{
							R: uint8(r),
							G: uint8(g),
							B: uint8(b),
							A: uint8(a),
						}
						newFrame.Set(x, y, newColor)
						doneEmbedding = true
					}
				}
			} // End x
		} // End y
		// Tweak the color palette to support the embedded LSB colors
		fmt.Println(previous, bitsIndex)
		colorPalette := TweakColorPalette(img.Palette, secret.Data[previous:bitsIndex])
		newFrame.Palette = colorPalette
		// Append the current frame since
		newGif.Image = append(newGif.Image, newFrame)
		if doneEmbedding {
			// Append the next frame and every frame after it and return
			newGif.Image = append(newGif.Image, file.Image[i+1:]...)
			return newGif, nil
		}
		previous = bitsIndex
	}
	return newGif, nil
}

func embedInColor(a byte, b uint8) uint8 {
	// 128 bit set indicates to zero out last 2 bits
	if a > 128 {
		b = b &^ 0x03 // zero out last 2 bits
		a = a & 3     // unset the 128 bit
	} else {
		b = b &^ 0x07 // zero out last 3 bits
	}
	c := b | a
	return c
}

// May make a comeback now that embed/extract works
//
// func embedIn16BitColor(a uint8, b uint32) uint16 {
// 	c := b | uint32(a)
// 	return uint16(c)
// 	// colorVal := strconv.FormatUint(uint64(b), 2)
// 	// colorBits := strings.Split(ZeroPadLeft(colorVal, 16), "")
// 	// // Least Significant Bits become bit pair of encoded msg
// 	// colorBits[len(colorBits)-2] = a[0]
// 	// colorBits[len(colorBits)-1] = a[1]

// 	// // Rejoin & reparse the new color values
// 	// colorStr := strings.Join(colorBits, "")
// 	// newColor, err := strconv.ParseUint(colorStr, 2, 64)
// 	// if err != nil {
// 	// 	return 0, err
// 	// }

// 	// return uint16(newColor), nil
// }
