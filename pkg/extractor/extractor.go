package extractor

import (
	"fmt"
	"image"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bshore/steggo/pkg/encoders"
	"github.com/bshore/steggo/pkg/process"
)

type Config struct {
	Target          io.ReadSeeker
	DestinationPath string
}

func Process(config *Config) error {
	var err error
	var message string
	_, format, err := image.Decode(config.Target)
	if err != nil {
		return fmt.Errorf("failed to decode target file: %v", err)
	}
	_, _ = config.Target.Seek(0, 0)

	switch format {
	case "png":
		message, err = ProcessPNG(config.Target)
		if err != nil {
			return fmt.Errorf("failed to process PNG: %v", err)
		}
	case "bmp":
		message, err = ProcessBMP(config.Target)
		if err != nil {
			return fmt.Errorf("failed to process BMP: %v", err)
		}
	case "gif":
		message, err = ProcessGif(config.Target)
		if err != nil {
			return fmt.Errorf("failed to process GIF: %v", err)
		}
	default:
		return fmt.Errorf("unsupported source file format: %v", format)
	}
	fmt.Println(message)
	if config.DestinationPath != "" {
		os.WriteFile(filepath.Join(config.DestinationPath, "message.txt"), []byte(message), fs.FileMode(os.O_WRONLY))
	}

	return nil
}

func DecodeMessage(header *process.Header, extracted []byte) (string, error) {
	msg := string(extracted)
	var err error
	if header.PreEncoding == "" {
		return string(extracted), nil
	}
	encStrings := strings.Split(header.PreEncoding, "/")
	preEncoders, warnings := encoders.FromIntStrSlice(encStrings)
	if warnings != "" {
		return "", fmt.Errorf("failed to determine pre-encoders (%v): %s", encStrings, warnings)
	}
	for _, enc := range preEncoders {
		switch enc {
		case encoders.R13:
			msg = encoders.Rot13(msg)
		case encoders.B16:
			msg, err = encoders.Decode16(msg)
			if err != nil {
				return "", fmt.Errorf("failed to decode b16 message: %v", err)
			}
		case encoders.B32:
			msg, err = encoders.Decode32(msg)
			if err != nil {
				return "", fmt.Errorf("failed to decode b32 message: %v", err)
			}
		case encoders.B64:
			msg, err = encoders.Decode64(msg)
			if err != nil {
				return "", fmt.Errorf("failed to decode b64 message: %v", err)
			}
		case encoders.B85:
			msg, err = encoders.Decode85(msg)
			if err != nil {
				return "", fmt.Errorf("failed to decode b85 message: %v", err)
			}
		default:
			return "", fmt.Errorf("attempted to decode unkown pre-encoding type: %d", enc)
		}
	}
	return msg, err
}
