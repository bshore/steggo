package process

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"reflect"
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
func EmbedMsgInGIF(data []byte, file *gif.GIF, frameFactor int) (*gif.GIF, error) {
	// early return if there's no way data will fit in GIF even at 100% color capacity
	if len(data) > GifMaxColor*len(file.Image) {
		return nil, fmt.Errorf("message will not fit in image")
	}
	// var skipIndex uint8
	// if file.BackgroundIndex != 0 {
	// 	skipIndex = file.BackgroundIndex
	// }
	// For each image frame
	for frameNum, frameImg := range file.Image {
		// Get the number of colors used in this frame's palette
		size := getColorPaletteSize(frameImg) - 1
		// if skipIndex != 0 {
		// 	size--
		// }
		// chunk data so that it will fit into this frame's palette
		var chunks [][]byte
		for i := 0; i < len(data); i += size {
			endAt := i + size
			if endAt > len(data) {
				endAt = len(data)
			}
			chunks = append(chunks, data[i:endAt])
		}

		// if we've run out of chunks to process, the rest of the GIF can remain unchanged
		// and we return early
		if frameNum > len(chunks)-1 {
			return file, nil
		}

		if frameNum%frameFactor == 0 {
			fmt.Println("modifying frame #", frameNum)
			// Replace the colors in the local color palette with the message data
			newPalette := ModifyGifFrameColorPalette(frameImg, chunks[frameNum])
			if frameNum != 0 {
				file.Image[frameNum].Palette = newPalette
				prevFrame := file.Image[frameNum-1]
				newPix := mergeFrames(prevFrame, frameImg)
				file.Image[frameNum].Pix = newPix
				file.Image[frameNum].Rect = prevFrame.Rect
				file.Image[frameNum].Stride = prevFrame.Stride
			}
		}

	}
	return file, nil
}

// getColorPaletteSize returns the maximum number of colors used in the given image frame
func getColorPaletteSize(frame *image.Paletted) int {
	colors := make(map[color.Color]struct{})
	for _, c := range frame.Palette {
		colors[c] = struct{}{}
	}
	return len(colors)
}

// TODO: animation is probably not working because the Pix slice on new frames is sometimes
//			 much shorter than the previous frame's Pix slice.
//
//			 If we know the stride width is 454, we can plot the Pix array into a 2d slice, and
//       using rectangle intersections only copy over the indices that overlap
/*
	Example:
				(0,0)-(454,390) (36,72)-(428,352)
				177060 109760
				(0,0)-(454,390) (107,136)-(119,146)
				177060 120
*/

// mergeFramePix uses the Pix array from the previous frame and current frame
// and merges them together, changing only the Pix indices that have changed
// from the previous frame.
func mergeFramePix(previousFrame *image.Paletted, currentFrame *image.Paletted) []uint8 {
	mergedPixels := make([]uint8, len(previousFrame.Pix))
	copy(mergedPixels, previousFrame.Pix)

	// for each Pix index and colorIdx pointing to the Palette from the prevous frame's Pix slice
	for i, colorIdx := range previousFrame.Pix {
		previousFrameColor := previousFrame.Palette[colorIdx]
		currentFrameColor := currentFrame.Palette[colorIdx]
		// check if the color from the palette of the previous has changed
		if !reflect.DeepEqual(previousFrameColor, currentFrameColor) {
			// this color has changed since the previous frame
			for j := range currentFrame.Pix {
				currentPaletteIdx := currentFrame.Pix[j]
				if currentFrame.Palette[currentPaletteIdx] == currentFrameColor {
					mergedPixels[i] = uint8(j)
					break
				}
			}
		} else {
			mergedPixels[i] = colorIdx
		}
	}
	return mergedPixels
}

func mergeFrames(prevFrame *image.Paletted, curFrame *image.Paletted) []uint8 {
	// Get the dimensions of the overlapping region between the two frames
	bounds := prevFrame.Bounds().Intersect(curFrame.Bounds())

	// Calculate the total number of pixels in the previous frame
	widthPrev, heightPrev := prevFrame.Bounds().Dx(), prevFrame.Bounds().Dy()
	numPixelsPrev := widthPrev * heightPrev

	// Create a new slice to store the merged pixels
	mergedPixels := make([]uint8, numPixelsPrev)

	// Copy the pixels from the previous frame to the merged pixels slice
	copy(mergedPixels, prevFrame.Pix)

	// Iterate over the overlapping region of the two frames
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the color index for the current pixel in each frame
			prevIndex := prevFrame.ColorIndexAt(x, y)
			curIndex := curFrame.ColorIndexAt(x, y)

			// If the color index has changed between frames, update the corresponding pixel in the merged pixels slice
			if prevIndex != curIndex {
				mergedPixels[(y-prevFrame.Rect.Min.Y)*widthPrev+(x-prevFrame.Rect.Min.X)] = curIndex
			}
		}
	}

	return mergedPixels
}

func mergeFrames2(previousFrame *image.Paletted, currentFrame *image.Paletted) []uint8 {
	// Find the intersecting rectangle of the two frames
	overlap := previousFrame.Bounds().Intersect(currentFrame.Bounds())

	// Create a 2D slice representing all the pixels of the previous frame
	pixels := make([][]uint8, previousFrame.Bounds().Dy())
	for y := range pixels {
		pixels[y] = previousFrame.Pix[y*previousFrame.Stride : (y+1)*previousFrame.Stride]
	}

	// Replace the indices of the 2D slice with the pixels from the current frame where the overlap is
	for y := overlap.Min.Y; y < overlap.Max.Y; y++ {
		for x := overlap.Min.X; x < overlap.Max.X; x++ {
			pixels[y][x] = currentFrame.Pix[currentFrame.PixOffset(x, y)]
		}
	}

	// Flatten the 2D slice into a 1D slice
	flattened := make([]uint8, 0, len(previousFrame.Pix))
	for _, row := range pixels {
		flattened = append(flattened, row...)
	}

	return flattened
}
