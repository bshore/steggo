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

// EmbedMsgInGIF takes the message string and embeds it into a GIF file
// frame by frame using Least Significant Bit(s)
func EmbedMsgInGIF(data []byte, file *gif.GIF) (*gif.GIF, error) {
	if len(data) > GifMaxColor*len(file.Image) {
		return nil, fmt.Errorf("message will not fit in image")
	}
	// Chunk the data into multiple 256 length []byte
	var chunks = make([][]byte, len(file.Image))
	var j int
	var done bool
	for i := 0; i < len(data); i += GifMaxColor {
		endAt := i + GifMaxColor
		if endAt > len(data) {
			endAt = len(data)
			done = true
		}
		chunks[j] = data[i:endAt]
		j++
		if done {
			break
		}
	}

	// For each image frame
	for frameNum, frameImg := range file.Image {
		newPallet := GetGifFrameColorPalette(frameImg, chunks[frameNum])
		file.Image[frameNum].Palette = newPallet
	}
	return file, nil
}
