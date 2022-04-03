package extractor

import (
	"fmt"
	"image/gif"
	"io"

	"github.com/bshore/steggo/pkg/process"
)

func ProcessGif(src io.Reader) (string, error) {
	loadedImage, err := gif.DecodeAll(src)
	if err != nil {
		return "", fmt.Errorf("error decoding GIF file: %v", err)
	}
	header, extracted, err := process.ExtractMsgFromGif(loadedImage)
	if err != nil {
		return "", fmt.Errorf("error extracting from GIF image: %v", err)
	}
	return DecodeMessage(header, extracted)
}
