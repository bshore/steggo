package embedder

import (
	"fmt"
	"image/color"
	"image/gif"
	"io"
	"os"

	"github.com/bshore/steggo/pkg/process"
)

func ProcessGIF(data []byte, dest string, src io.Reader) error {
	loaded, err := gif.DecodeAll(src)
	if err != nil {
		return fmt.Errorf("error decoding GIF file: %v", err)
	}

	// Convert global color table to local color tables if present
	if loaded.Config.ColorModel != nil {
		globalColorTable := loaded.Config.ColorModel.(color.Palette)
		// Place the global color table as the local color table of all frames
		for i := range loaded.Image {
			if loaded.Image[i].Palette == nil {
				paletteCopy := make(color.Palette, len(globalColorTable))
				copy(paletteCopy, globalColorTable)
				loaded.Image[i].Palette = paletteCopy
			}
		}
		// Clear the global color table
		loaded.Config.ColorModel = nil
	}

	embedded, err := process.EmbedMsgInGIF(data, loaded)
	if err != nil {
		return fmt.Errorf("error embedding message in file: %v", err)
	}

	newFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer newFile.Close()

	err = gif.EncodeAll(newFile, embedded)
	if err != nil {
		return fmt.Errorf("error encoding new GIF image: %v", err)
	}
	return nil
}
