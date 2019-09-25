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
	var newR, newG, newB uint16
	bounds := file.Bounds()
	pixels := bounds.Max.X * bounds.Max.Y
	if !(secret.Size < pixels) {
		return nil, fmt.Errorf("Secret message won't fit in image: %v LSB's to embed, %v pixels available", secret.Size, pixels)
	}
	newFile := image.NewNRGBA64(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
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

// EmbedMsgInGIF takes the message string and embeds it into a GIF file
// frame by frame using Least Significant Bit(s)
func EmbedMsgInGIF(secret *Secret, file *gif.GIF) (*gif.GIF, error) {
	var bitsIndex int
	var doneEmbedding bool
	var newR, newG, newB uint8
	var newFrame *image.Paletted
	var colorPalette []color.Color
	// // Color table only allows for 256 color combinations, multiply by number of frames for available pixels.
	// pixels := 256 * len(file.Image)
	// // If the secret's size is not under the amount of available pixels, we can't embed.
	// if !(secret.Size/3 < pixels) {
	// 	return nil, fmt.Errorf("Secret message won't fit in image: %v LSB's to embed, %v pixels available", secret.Size/3, pixels)
	// }
	newGif := &gif.GIF{
		Image:           []*image.Paletted{},
		Delay:           file.Delay,
		LoopCount:       file.LoopCount,
		Disposal:        file.Disposal,
		Config:          file.Config,
		BackgroundIndex: file.BackgroundIndex,
	}
	// For each image frame
	for frameNum, frameImg := range file.Image {
		// The image rectangle bounds
		bounds := frameImg.Bounds()
		// An empty frame with the same size as the source GIF and an empty color palette
		newFrame = image.NewPaletted(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), nil)
		// colorPalette = CreateColorPalette(frameImg.Palette, secret.Message)
		colorPalette = GetGifFrameColorPalette(frameImg, secret.Message, secret.Data)
		if doneEmbedding {
			newFrame.Palette = frameImg.Palette
		} else {
			newFrame.Palette = colorPalette
		}
		// For each vertical row
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			// For each pixel in each row
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := frameImg.At(x, y).RGBA()
				// if bitsIndex >= GifMaxPerFrame*(frameNum+1) {
				// 	newFrame.Set(x, y, frameImg.At(x, y))
				// 	continue
				// }
				// If the iteration is still under the length of message bits
				if bitsIndex < secret.Size {
					newR = embedInColor(secret.Data[bitsIndex], uint8(r))
					// Check if next msg byte to embed and if byte will fit
					if bitsIndex+1 < secret.Size && (bitsIndex+1)%GifMaxPerFrame != 0 {
						newG = embedInColor(secret.Data[bitsIndex+1], uint8(g))
						// Check if next msg byte to embed and if byte will fit
						if bitsIndex+2 < secret.Size && (bitsIndex+1)%GifMaxPerFrame != 0 {
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
			} // End x
		} // End y
		// If bitsIndex greater than GifMaxPerFrame * current frame, +1 for next frame
		if bitsIndex == GifMaxPerFrame*(frameNum+1) {
			bitsIndex++
		}
		// Append the current frame since
		newGif.Image = append(newGif.Image, newFrame)
		if doneEmbedding {
			// Append the next frame and every frame after it and return
			newGif.Image = append(newGif.Image, file.Image[frameNum+1:]...)
			return newGif, nil
		}
	}
	return newGif, nil
}
