package process

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"sort"
)

/*
	This file contains generic struct types and helper functions
*/

// GifMaxColor is the maximum amount of colors that are supported by a frame's
// Local Color Table.
const GifMaxColor int = 256

// GifMaxPerFrame is the sum of RGB pixels for which embedding can occur per frame.
// Each color is made up of 3 bytes, Local Color Table has a max of 256 colors:
// 3 * 256 = 768
const GifMaxPerFrame int = 768

// Flags holds the types of flags allowed by the script
type Flags struct {
	SrcFile     string
	OutputDir   string
	Text        string
	Stdin       bool
	MessageFile string
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

/*
	LSB's
	64 of each for red: (4x64 = 256)
		00, 01, 10, 11
	32 of each for green and blue: (8x32 = 256)
		000, 001, 010, 011, 100, 101, 110, 111

	"Alphabetize" the frame's color palette beforehand, regardless of message ?
	Find a way to make sure enough colors are represented ?

	^ In the order of the "secret" message, register alphabet letters uniquely, and then
	register the remainder

	^ Make a map of map[byte]color, and then when drawing, use the key to draw the color.


	Maybe an algorithm that finds most common/ least common characters n stuff
	and reserves colors for those characters or something.

	Maybe GIF does a mapping instead of just pixels ?
	First character stores pointer to next character x/y coordinate ?
	First 256 characters store mappings to get words/paragraphs ?
*/

// chars := make(map[string]int)
// 	for _, c := range msg {
// 		if _, ok := chars[string(c)]; !ok {
// 			chars[string(c)] = 0
// 		}
// 		chars[string(c)]++
// 	}
// 	fmt.Println("Unique Characters: ", len(chars))
// 	for ch, ct := range chars {
// 		fmt.Printf("Character: '%v'\tCount: %v\n", ch, ct)
// 	}
// 	y := bounds.Min.Y
// 	x := bounds.Min.X
// 	// Make a diagonal pass from 0,0 to Ymax, Xmax to get a color palette
// 	for y < bounds.Max.Y && x < bounds.Max.X {}

// Alphabet is every printable character split in a 2-3-3
// First index of each slice has 128 bit set so we can identify
// later on when to zero out which least significant bits. Bit 128
// signifies to zero out last 2, an unset 128 means zero out last 3.
// Shift bb------ to ------bb and get last 2 bits
// Shift --bbb--- to -----bbb and get last 3 bits
// Just  -----bbb get last 3 bits
var Alphabet = [][]byte{
	[]byte{(0x20 >> 6) & 131, (0x20 >> 3) & 7, 0x20 & 7}, // (space)
	[]byte{(0x21 >> 6) & 131, (0x21 >> 3) & 7, 0x21 & 7}, // !
	[]byte{(0x22 >> 6) & 131, (0x22 >> 3) & 7, 0x22 & 7}, // "
	[]byte{(0x23 >> 6) & 131, (0x23 >> 3) & 7, 0x23 & 7}, // #
	[]byte{(0x24 >> 6) & 131, (0x24 >> 3) & 7, 0x24 & 7}, // $
	[]byte{(0x25 >> 6) & 131, (0x25 >> 3) & 7, 0x25 & 7}, // %
	[]byte{(0x26 >> 6) & 131, (0x26 >> 3) & 7, 0x26 & 7}, // &
	[]byte{(0x27 >> 6) & 131, (0x27 >> 3) & 7, 0x27 & 7}, // '
	[]byte{(0x28 >> 6) & 131, (0x28 >> 3) & 7, 0x28 & 7}, // (
	[]byte{(0x29 >> 6) & 131, (0x29 >> 3) & 7, 0x29 & 7}, // )
	[]byte{(0x2A >> 6) & 131, (0x2A >> 3) & 7, 0x2A & 7}, // *
	[]byte{(0x2B >> 6) & 131, (0x2B >> 3) & 7, 0x2B & 7}, // +
	[]byte{(0x2C >> 6) & 131, (0x2C >> 3) & 7, 0x2C & 7}, // ,
	[]byte{(0x2D >> 6) & 131, (0x2D >> 3) & 7, 0x2D & 7}, // -
	[]byte{(0x2E >> 6) & 131, (0x2E >> 3) & 7, 0x2E & 7}, // .
	[]byte{(0x2F >> 6) & 131, (0x2F >> 3) & 7, 0x2F & 7}, // /
	[]byte{(0x30 >> 6) & 131, (0x30 >> 3) & 7, 0x30 & 7}, // 0
	[]byte{(0x31 >> 6) & 131, (0x31 >> 3) & 7, 0x31 & 7}, // 1
	[]byte{(0x32 >> 6) & 131, (0x32 >> 3) & 7, 0x32 & 7}, // 2
	[]byte{(0x33 >> 6) & 131, (0x33 >> 3) & 7, 0x33 & 7}, // 3
	[]byte{(0x34 >> 6) & 131, (0x34 >> 3) & 7, 0x34 & 7}, // 4
	[]byte{(0x35 >> 6) & 131, (0x35 >> 3) & 7, 0x35 & 7}, // 5
	[]byte{(0x36 >> 6) & 131, (0x36 >> 3) & 7, 0x36 & 7}, // 6
	[]byte{(0x37 >> 6) & 131, (0x37 >> 3) & 7, 0x37 & 7}, // 7
	[]byte{(0x38 >> 6) & 131, (0x38 >> 3) & 7, 0x38 & 7}, // 8
	[]byte{(0x39 >> 6) & 131, (0x39 >> 3) & 7, 0x39 & 7}, // 9
	[]byte{(0x3A >> 6) & 131, (0x3A >> 3) & 7, 0x3A & 7}, // :
	[]byte{(0x3B >> 6) & 131, (0x3B >> 3) & 7, 0x3B & 7}, // ;
	[]byte{(0x3C >> 6) & 131, (0x3C >> 3) & 7, 0x3C & 7}, // <
	[]byte{(0x3D >> 6) & 131, (0x3D >> 3) & 7, 0x3D & 7}, // =
	[]byte{(0x3E >> 6) & 131, (0x3E >> 3) & 7, 0x3E & 7}, // >
	[]byte{(0x3F >> 6) & 131, (0x3F >> 3) & 7, 0x3F & 7}, // ?
	[]byte{(0x40 >> 6) & 131, (0x40 >> 3) & 7, 0x40 & 7}, // @
	[]byte{(0x41 >> 6) & 131, (0x41 >> 3) & 7, 0x41 & 7}, // A
	[]byte{(0x42 >> 6) & 131, (0x42 >> 3) & 7, 0x42 & 7}, // B
	[]byte{(0x43 >> 6) & 131, (0x43 >> 3) & 7, 0x43 & 7}, // C
	[]byte{(0x44 >> 6) & 131, (0x44 >> 3) & 7, 0x44 & 7}, // D
	[]byte{(0x45 >> 6) & 131, (0x45 >> 3) & 7, 0x45 & 7}, // E
	[]byte{(0x46 >> 6) & 131, (0x46 >> 3) & 7, 0x46 & 7}, // F
	[]byte{(0x47 >> 6) & 131, (0x47 >> 3) & 7, 0x47 & 7}, // G
	[]byte{(0x48 >> 6) & 131, (0x48 >> 3) & 7, 0x48 & 7}, // H
	[]byte{(0x49 >> 6) & 131, (0x49 >> 3) & 7, 0x49 & 7}, // I
	[]byte{(0x4A >> 6) & 131, (0x4A >> 3) & 7, 0x4A & 7}, // J
	[]byte{(0x4B >> 6) & 131, (0x4B >> 3) & 7, 0x4B & 7}, // K
	[]byte{(0x4C >> 6) & 131, (0x4C >> 3) & 7, 0x4C & 7}, // L
	[]byte{(0x4D >> 6) & 131, (0x4D >> 3) & 7, 0x4D & 7}, // M
	[]byte{(0x4E >> 6) & 131, (0x4E >> 3) & 7, 0x4E & 7}, // N
	[]byte{(0x4F >> 6) & 131, (0x4F >> 3) & 7, 0x4F & 7}, // O
	[]byte{(0x50 >> 6) & 131, (0x50 >> 3) & 7, 0x50 & 7}, // P
	[]byte{(0x51 >> 6) & 131, (0x51 >> 3) & 7, 0x51 & 7}, // Q
	[]byte{(0x52 >> 6) & 131, (0x52 >> 3) & 7, 0x52 & 7}, // R
	[]byte{(0x53 >> 6) & 131, (0x53 >> 3) & 7, 0x53 & 7}, // S
	[]byte{(0x54 >> 6) & 131, (0x54 >> 3) & 7, 0x54 & 7}, // T
	[]byte{(0x55 >> 6) & 131, (0x55 >> 3) & 7, 0x55 & 7}, // U
	[]byte{(0x56 >> 6) & 131, (0x56 >> 3) & 7, 0x56 & 7}, // V
	[]byte{(0x57 >> 6) & 131, (0x57 >> 3) & 7, 0x57 & 7}, // W
	[]byte{(0x58 >> 6) & 131, (0x58 >> 3) & 7, 0x58 & 7}, // X
	[]byte{(0x59 >> 6) & 131, (0x59 >> 3) & 7, 0x59 & 7}, // Y
	[]byte{(0x5A >> 6) & 131, (0x5A >> 3) & 7, 0x5A & 7}, // Z
	[]byte{(0x5B >> 6) & 131, (0x5B >> 3) & 7, 0x5B & 7}, // [
	[]byte{(0x5C >> 6) & 131, (0x5C >> 3) & 7, 0x5C & 7}, // \
	[]byte{(0x5D >> 6) & 131, (0x5D >> 3) & 7, 0x5D & 7}, // ]
	[]byte{(0x5E >> 6) & 131, (0x5E >> 3) & 7, 0x5E & 7}, // ^
	[]byte{(0x5F >> 6) & 131, (0x5F >> 3) & 7, 0x5F & 7}, // _
	[]byte{(0x60 >> 6) & 131, (0x60 >> 3) & 7, 0x60 & 7}, // `
	[]byte{(0x61 >> 6) & 131, (0x61 >> 3) & 7, 0x61 & 7}, // a
	[]byte{(0x62 >> 6) & 131, (0x62 >> 3) & 7, 0x62 & 7}, // b
	[]byte{(0x63 >> 6) & 131, (0x63 >> 3) & 7, 0x63 & 7}, // c
	[]byte{(0x64 >> 6) & 131, (0x64 >> 3) & 7, 0x64 & 7}, // d
	[]byte{(0x65 >> 6) & 131, (0x65 >> 3) & 7, 0x65 & 7}, // e
	[]byte{(0x66 >> 6) & 131, (0x66 >> 3) & 7, 0x66 & 7}, // f
	[]byte{(0x67 >> 6) & 131, (0x67 >> 3) & 7, 0x67 & 7}, // g
	[]byte{(0x68 >> 6) & 131, (0x68 >> 3) & 7, 0x68 & 7}, // h
	[]byte{(0x69 >> 6) & 131, (0x69 >> 3) & 7, 0x69 & 7}, // i
	[]byte{(0x6A >> 6) & 131, (0x6A >> 3) & 7, 0x6A & 7}, // j
	[]byte{(0x6B >> 6) & 131, (0x6B >> 3) & 7, 0x6B & 7}, // k
	[]byte{(0x6C >> 6) & 131, (0x6C >> 3) & 7, 0x6C & 7}, // l
	[]byte{(0x6D >> 6) & 131, (0x6D >> 3) & 7, 0x6D & 7}, // m
	[]byte{(0x6E >> 6) & 131, (0x6E >> 3) & 7, 0x6E & 7}, // n
	[]byte{(0x6F >> 6) & 131, (0x6F >> 3) & 7, 0x6F & 7}, // o
	[]byte{(0x70 >> 6) & 131, (0x70 >> 3) & 7, 0x70 & 7}, // p
	[]byte{(0x71 >> 6) & 131, (0x71 >> 3) & 7, 0x71 & 7}, // q
	[]byte{(0x72 >> 6) & 131, (0x72 >> 3) & 7, 0x72 & 7}, // r
	[]byte{(0x73 >> 6) & 131, (0x73 >> 3) & 7, 0x73 & 7}, // s
	[]byte{(0x74 >> 6) & 131, (0x74 >> 3) & 7, 0x74 & 7}, // t
	[]byte{(0x75 >> 6) & 131, (0x75 >> 3) & 7, 0x75 & 7}, // u
	[]byte{(0x76 >> 6) & 131, (0x76 >> 3) & 7, 0x76 & 7}, // v
	[]byte{(0x77 >> 6) & 131, (0x77 >> 3) & 7, 0x77 & 7}, // w
	[]byte{(0x78 >> 6) & 131, (0x78 >> 3) & 7, 0x78 & 7}, // x
	[]byte{(0x79 >> 6) & 131, (0x79 >> 3) & 7, 0x79 & 7}, // y
	[]byte{(0x7A >> 6) & 131, (0x7A >> 3) & 7, 0x7A & 7}, // z
	[]byte{(0x7B >> 6) & 131, (0x7B >> 3) & 7, 0x7B & 7}, // {
	[]byte{(0x7C >> 6) & 131, (0x7C >> 3) & 7, 0x7C & 7}, // |
	[]byte{(0x7D >> 6) & 131, (0x7D >> 3) & 7, 0x7D & 7}, // }
	[]byte{(0x7E >> 6) & 131, (0x7E >> 3) & 7, 0x7E & 7}, // ~
}

type charCount struct {
	Char  byte
	Count int
}

func (c charCount) lsbSplit() (byte, byte, byte) {
	return (c.Char >> 6) & 131, (c.Char >> 3) & 7, c.Char & 7
}

// CreateColorPalette is a piece of crap
// func CreateColorPalette(old []color.Color, msg []byte) []color.Color {
// 	var new []color.Color
// 	var chars []charCount
// 	var i int
// 	var col color.Color
// 	var charIndex = 0
// 	charMap := make(map[byte]int)
// 	for _, c := range msg {
// 		if _, ok := charMap[c]; !ok {
// 			charMap[c] = 0
// 		}
// 		charMap[c]++
// 	}
// 	for k, v := range charMap {
// 		chars = append(chars, charCount{Char: k, Count: v})
// 	}
// 	sort.Slice(chars, func(i, j int) bool {
// 		return chars[i].Count > chars[j].Count
// 	})
// 	fmt.Println(chars)
// 	// Make sure every printable Alphabet character is represented
// 	for len(new) < 256 {
// 		i = 0
// 		for i, col = range old {
// 			if len(new) == 256 {
// 				break
// 			}
// 			r, g, b, a := col.RGBA()
// 			if i < len(Alphabet) {
// 				newR := embedInColor(Alphabet[i][0], uint8(r))
// 				newG := embedInColor(Alphabet[i][1], uint8(g))
// 				newB := embedInColor(Alphabet[i][2], uint8(b))
// 				new = append(new, color.RGBA{
// 					R: newR,
// 					G: newG,
// 					B: newB,
// 					A: uint8(a),
// 				})
// 			} else {
// 				// Now designate more colors for more used characters
// 				lsbR, lsbG, lsbB := chars[charIndex].lsbSplit()
// 				newR := embedInColor(lsbR, uint8(r))
// 				newG := embedInColor(lsbG, uint8(g))
// 				newB := embedInColor(lsbB, uint8(b))
// 				new = append(new, color.RGBA{
// 					R: newR,
// 					G: newG,
// 					B: newB,
// 					A: uint8(a),
// 				})
// 				if chars[charIndex].Count < 10 {
// 					charIndex = 0
// 				}
// 				charIndex++
// 			}
// 		}
// 	}
// 	return new
// }

// GetGifFrameColorPalette gathers a Gif Frame's Color Palette
func GetGifFrameColorPalette(img *image.Paletted, msg []byte, data []byte) []color.Color {
	var bitsIndex int
	var newR, newG, newB uint8
	var colorPalette []color.Color
	paletteMap := make(map[color.Color]struct{})
	bounds := img.Bounds()
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if bitsIndex < len(data) && bitsIndex < GifMaxPerFrame {
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
				newColor := &color.RGBA{
					R: newR,
					G: newG,
					B: newB,
					A: uint8(a),
				}
				if _, ok := paletteMap[newColor]; !ok {
					paletteMap[newColor] = struct{}{}
				}
			}
			bitsIndex += 3
		}
	}
	var j int
	for len(paletteMap) < 256 && j < len(img.Palette) {
		if _, ok := paletteMap[img.Palette[j]]; !ok {
			paletteMap[img.Palette[j]] = struct{}{}
		}
		j++
	}
	for c := range paletteMap {
		colorPalette = append(colorPalette, c)
	}
	sort.SliceStable(colorPalette, func(i, j int) bool {
		r1, g1, b1, _ := colorPalette[i].RGBA()
		r2, g2, b2, _ := colorPalette[j].RGBA()
		return (r1 + g1 + b1) < (r2 + g2 + b2)
	})
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
