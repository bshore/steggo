package embedder

import (
	"fmt"
	"image"
	"io"
	"lsb_encoder/pkg/encoders"
	"lsb_encoder/pkg/process"
	"path/filepath"
)

type Config struct {
	Input           string
	SrcType         string
	Target          io.ReadSeeker
	DestinationPath string
	PreEncoding     []encoders.EncType
}

func Process(config *Config) error {
	processedInput := encoders.ApplyPreEncoding(config.Input, config.PreEncoding)

	_, format, err := image.Decode(config.Target)
	if err != nil {
		return fmt.Errorf("failed to decode target file: %v", err)
	}
	_, _ = config.Target.Seek(0, 0)

	dest := formatDestination(config.DestinationPath, format)
	header := process.NewHeaderBytes(processedInput, config.SrcType, config.PreEncoding)
	data := process.FinalizeMessage(header, processedInput)

	switch format {
	case "png":
		return ProcessPNG(data, dest, config.Target)
	case "jpeg":
		return ProcessJPEG(data, dest, config.Target)
	case "bmp":
		return ProcessBMP(data, dest, config.Target)
	// case "gif":
	// 	return ProcessGIF(data, dest, config.Target)
	default:
		return fmt.Errorf("unsupported source file format: %v", format)
	}
}

// formatDestination returns output.{format} except for jpeg, which returns output_jpeg.png
//
//  The reason for outputting a png for jpeg input is due to jpeg's native compression, we
//  don't want to output jpeg since the simple act of saving a jpeg risks destroying the
//  embedded message.
func formatDestination(path, format string) string {
	if format == "jpeg" {
		return filepath.Join(path, fmt.Sprintf("output_%s.png", format))
	}
	return filepath.Join(path, fmt.Sprintf("output.%s", format))
}
