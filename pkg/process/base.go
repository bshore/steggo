package process

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
)

/*
	This file contains generic struct types and helper functions
*/

// GifMaxColor is the maximum amount of colors that are supported by a frame's
// Local Color Table.
const GifMaxColor int = 256

// Flags holds the types of flags allowed by the script
type Flags struct {
	SrcFile     string
	OutputDir   string
	Text        string
	Stdin       bool
	MessageFile string
	BitOpt      int
	Decode      bool
	Rot13       bool
	Base16      bool
	Base32      bool
	Base64      bool
	Base85      bool
	Complex     string
}

// ZeroPadLeft left pads a string with zeros until the string is
// <length> characters long
func ZeroPadLeft(str string, length int) string {
	if len(str) == length {
		return str
	}
	for {
		str = "0" + str
		if len(str) == length {
			return str
		}
	}
}

// WriteFile writes the new file to the output directory
func WriteFile(data []byte, out, ext string) error {
	err := os.MkdirAll(out, 0777)
	if err != nil {
		return err
	}
	newFile, err := os.Create(filepath.Join(out, "output"+ext))
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = newFile.Write(data)
	return err
}

// GetGifFrameColorPalette gathers a Gif Frame's Color Palette
func GetGifFrameColorPalette(img *image.Paletted, data []byte) []color.Color {
	if len(data) == 0 {
		return img.Palette
	}
	var bitsIndex int
	var newR, newG, newB uint8
	var colorPalette []color.Color
	for _, paletteColor := range img.Palette {
		r, g, b, a := paletteColor.RGBA()
		if bitsIndex < len(data) && bitsIndex < GifMaxColor {
			newR = embedInColor(data[bitsIndex], uint8(r))
			if bitsIndex+1 < len(data) {
				newG = embedInColor(data[bitsIndex+1], uint8(g))
				if bitsIndex+2 < len(data) {
					newB = embedInColor(data[bitsIndex+2], uint8(b))
				} else {
					newB = uint8(b)
				}
			} else {
				newG = uint8(g)
			}
			colorPalette = append(colorPalette, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a)})
		}
		bitsIndex += 3
	}
	return colorPalette
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

func embedIn16BitColor(a uint8, b uint32) uint16 {
	if a > 128 {
		b = b &^ 0x03 // zero out last 2 bits
		a = a & 3     // unset the 128 bit
	} else {
		b = b &^ 0x07 // zero out the last 3 bits
	}
	c := b | uint32(a)
	return uint16(c)
}

func extractFromColor(r, g, b uint8) byte {
	// Get last bits of each color to reconstruct a message byte
	rBits := r & 3
	gBits := g & 7
	bBits := b & 7

	var newByte uint8
	// Assign color bits and shift left for each color pixel
	newByte = newByte | (rBits & 3) // ------bb
	newByte = newByte << 3          // ---bb---
	newByte = newByte | (gBits & 7) // -----bbb
	newByte = newByte << 3          // bbbbb---
	newByte = newByte | (bBits & 7) // bbbbbbbb
	return newByte
}

// func rebuildFromBits(b []uint8) byte {
// 	var newByte uint8
// 	newByte = newByte | b[0]
// 	newByte = newByte << 1
// 	newByte = newByte | b[1]
// 	newByte = newByte << 1
// 	newByte = newByte | b[2]
// 	newByte = newByte << 1
// 	newByte = newByte | b[3]
// 	newByte = newByte << 1
// 	newByte = newByte | b[4]
// 	newByte = newByte << 1
// 	newByte = newByte | b[5]
// 	newByte = newByte << 1
// 	newByte = newByte | b[6]
// 	newByte = newByte << 1
// 	newByte = newByte | b[7]
// 	return newByte
// }
