package embedder

import (
	"fmt"
	"image"
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

	if !delayIsConsistent(loaded.Delay) {
		return fmt.Errorf("error: cannot process gif with inconsistent frame delay")
	}

	frameFactor := getFrameFactor(loaded.Delay[0])
	disposal := loaded.Disposal[0]
	newDelay := loaded.Delay[0] / frameFactor
	fmt.Println("delay is", loaded.Delay[0])
	fmt.Println("frame factor is", frameFactor)
	fmt.Println("new delay is", newDelay)
	fmt.Println("len og frames", len(loaded.Image))

	dupedFrames := make([]image.Paletted, len(loaded.Image)*frameFactor)
	fmt.Println("len duped frames", len(dupedFrames))
	for i, img := range loaded.Image {
		dupedFrames[i] = *img
		for j := 1; j < frameFactor; j++ {
			dupedFrames[i+j] = *img
		}
	}

	loaded.Disposal = make([]byte, len(dupedFrames))
	for i := range loaded.Disposal {
		loaded.Disposal[i] = disposal
	}

	loaded.Delay = make([]int, len(dupedFrames))
	for i := range loaded.Delay {
		loaded.Delay[i] = newDelay
	}

	loaded.Image = make([]*image.Paletted, len(dupedFrames))
	for i := range loaded.Image {
		loaded.Image[i] = &dupedFrames[i]
	}

	if loaded.Config.ColorModel != nil {
		// This gif uses the global color table, we need to pull it out and make
		// each frame use it as a starting point for its local color table.
		globalColorTable := loaded.Config.ColorModel.(color.Palette)
		// Place the global color table as the local color table of
		// all frames
		for i := range loaded.Image {
			loaded.Image[i].Palette = globalColorTable
		}
		// and clear the global color table
		loaded.Config.ColorModel = nil
	}
	embedded, err := process.EmbedMsgInGIF(data, loaded, frameFactor)
	if err != nil {
		return fmt.Errorf("error embedding message in file: %v", err)
	}
	newFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer newFile.Close()

	embedded.Config.ColorModel = nil // clear the global color table
	err = gif.EncodeAll(newFile, embedded)
	if err != nil {
		return fmt.Errorf("error encoding new GIF image: %v", err)
	}
	return nil
}

func delayIsConsistent(delays []int) bool {
	if len(delays) == 0 {
		return true
	}
	first := delays[0]
	for i := range delays {
		if delays[i] != first {
			return false
		}
	}
	return true
}

func getFrameFactor(delay int) int {
	for i := 2; i <= delay; i++ {
		if delay%i == 0 {
			return i
		}
	}
	return 0
}
