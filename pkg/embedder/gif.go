package embedder

import (
	"fmt"
	"image/color"
	"image/gif"
	"io"
	"os"

	"github.com/bshore/steggo/pkg/process"
)

/* TODO: instead of manipulating the GIF in place
insert an altered 0th frame with metadata about
what a skip factor would be

skip factor may be determined by the frame delay

ex: if frame delay is 5, every 5th frame is an LSB
altered frame, and the delay is adjusted to 1
with frames 1-4 being copied original frames

for ex, M = modified frame, c = copied frame

MccccMccccMccccM

the copied frame is the original, and with a shortened
delay keeping the full animation time the same, hopefully
the eye cannot notice

FRICK, this isn't working
TODO: try again
*/

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
				loaded.Image[i].Palette = globalColorTable
			}
		}
		// Clear the global color table
		loaded.Config.ColorModel = nil
	}

	// Embed the message using unused palette slots and disposal methods
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
