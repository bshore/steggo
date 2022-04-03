package embedder

import (
	"fmt"
	"image"
	"io"
	"path/filepath"

	"github.com/bshore/steggo/pkg/encoders"
	"github.com/bshore/steggo/pkg/process"
)

type Config struct {
	Input           string
	SrcType         string
	SrcFilename     string
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

	dest := formatDestination(config.SrcFilename, config.DestinationPath, format)
	header := process.NewHeaderBytes(processedInput, config.SrcType, config.PreEncoding)
	data := process.FinalizeMessage(header, processedInput)

	switch format {
	case "png":
		return ProcessPNG(data, dest, config.Target)
	case "jpeg":
		return ProcessJPEG(data, dest, config.Target)
	case "bmp":
		return ProcessBMP(data, dest, config.Target)
	case "gif":
		return ProcessGIF(data, dest, config.Target)
	default:
		return fmt.Errorf("unsupported source file format: %v", format)
	}
}

// formatDestination returns output.{format} except for jpeg and bmp, which returns output_<format>.png
//
//  The reason for outputting a .png for jpeg input is due to jpeg's native compression, we
//  don't want to output jpeg since the simple act of saving a jpeg risks destroying the
//  embedded message.
//
//  The reason for outputting a .png for bmp has to do with bmp only supporting 256 colors, so to avoid
//  embedding a message that can never be retrieved, we save the output as a .png
func formatDestination(srcFilename, path, format string) string {
	if format == "jpeg" || format == "jpg" || format == "bmp" {
		return filepath.Join(path, fmt.Sprintf("%s_%s_output.png", srcFilename, format))
	}
	return filepath.Join(path, fmt.Sprintf("%s_output.%s", srcFilename, format))
}
