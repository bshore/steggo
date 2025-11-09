package embedder

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/bshore/steggo/pkg/process"
)

func ProcessJPEG(data []byte, dest string, src io.Reader) error {
	loadedImage, err := jpeg.Decode(src)
	if err != nil {
		return fmt.Errorf("error decoding JPEG file: %v", err)
	}
	embedded, err := process.EmbedMsgInImage(data, loadedImage)
	if err != nil {
		return fmt.Errorf("error embedding message in file: %v", err)
	}
	newFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer newFile.Close()

	err = png.Encode(newFile, embedded)
	if err != nil {
		return fmt.Errorf("error encoding new JPEG image: %v", err)
	}
	return nil
}
