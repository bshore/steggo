package extractor

import (
	"fmt"
	"io"
	"lsb_encoder/pkg/process"

	"golang.org/x/image/bmp"
)

func ProcessBMP(src io.Reader) (string, error) {
	loadedImage, err := bmp.Decode(src)
	if err != nil {
		return "", fmt.Errorf("error decoding BMP file: %v", err)
	}
	header, extracted, err := process.ExtractMsgFromImage(loadedImage)
	if err != nil {
		return "", fmt.Errorf("error extracting from image: %v", err)
	}
	return DecodeMessage(header, extracted)
}
