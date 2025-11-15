package process

/*
	This file contains generic struct types and helper functions
*/

// GifMaxColor is the maximum amount of colors that are supported by a frame's
// Local Color Table.
const GifMaxColor int = 256

// GifFrameCapacity holds information about embedding capacity for a single frame
type GifFrameCapacity struct {
	FrameIndex       int
	UnusedIndices    []int
	DisposalMethod   byte
	CanModifyPalette bool
	Capacity         int // in bits (each color can hold 2+3+3=8 bits)
}

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

func embedInColor(a uint8, color uint8) uint8 {
	// 128 bit set indicates to zero out last 2 bits
	if a > 128 {
		color = color &^ 0x03 // zero out last 2 bits
		a = a & 3             // unset the 128 bit
	} else {
		color = color &^ 0x07 // zero out the last 3 bits
	}
	c := color | a
	return c
}

func embedIn16BitColor(a uint8, color uint32) uint16 {
	if a > 128 {
		color = color &^ 0x03 // zero out last 2 bits
		a = a & 3             // unset the 128 bit
	} else {
		color = color &^ 0x07 // zero out the last 3 bits
	}
	c := color | uint32(a)
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
