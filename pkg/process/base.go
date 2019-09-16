package process

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
)

/*
	This file contains generic struct types and helper functions
*/

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
	newFile, err := os.Create(filepath.Join(out, "output."+ext))
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = newFile.Write(data)
	return err
}

/*
	GIF embedding is totally screwed up.
	The Local Color Table needs to be tweaked on the fly somehow to
	allow for the LSB embedding to actually work, because on write/save
	the color table basicallly "rounds" whatever we LSB encoded to it's
	closest color value in the color table, blegh
*/

// TweakColorPalette adjusts a frame's color palette to support the adjusted LSB colors
func TweakColorPalette(oldPalette []color.Color, data []byte) []color.Color {
	fmt.Printf("Palette: %v\tData: %v\n", len(oldPalette), len(data))
	newPalette := []color.Color{}
	for i, col := range oldPalette {
		r, g, b, a := col.RGBA()
		newR := embedInColor(data[i], uint8(r))
		newG := embedInColor(data[i+1], uint8(g))
		newB := embedInColor(data[i+2], uint8(b))
		newPalette = append(newPalette, color.RGBA{
			R: newR, G: newG, B: newB, A: uint8(a),
		})
	}
	return newPalette
}
