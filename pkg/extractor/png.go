package extractor

import (
	"fmt"
	"image/png"
	"io"
	"lsb_encoder/pkg/process"
)

func ProcessPNG(src io.Reader) (string, error) {
	loadedImage, err := png.Decode(src)
	if err != nil {
		return "", fmt.Errorf("error decoding PNG file: %v", err)
	}
	header, extracted, err := process.ExtractMsgFromImage(loadedImage)
	if err != nil {
		return "", fmt.Errorf("error extracting from image: %v", err)
	}
	return DecodeMessage(header, extracted)
}
