package embed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bshore/steggo/pkg/embedder"
	"github.com/bshore/steggo/pkg/encoders"
	"github.com/bshore/steggo/pkg/utils"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "embed",
	Short: "",
	Long:  "",
	RunE:  embedCmdFn,
}

var (
	targetFile      string
	destinationPath string
	inputStr        string
	preEncoding     []string
)

func InitCmd() {
	Cmd.PersistentFlags().StringVar(&targetFile, "target", "", "The path to the image file being targeted for embedding")
	Cmd.PersistentFlags().StringVar(&destinationPath, "dest", ".", "The destination path to output the target file after embedding")
	Cmd.PersistentFlags().StringVar(&inputStr, "input", "", "The input path or message to embed into the target file")
	Cmd.PersistentFlags().StringSliceVar(&preEncoding, "pre-encoding", []string{}, "A comma separated list of pre-encoders to apply before embedding (r13, b16, b32, b64, b85), 5 maximum")
}

func embedCmdFn(command *cobra.Command, args []string) (err error) {
	input, srcType, err := getInputString(inputStr)
	if err != nil {
		return err
	}

	target, err := os.Open(targetFile)
	if err != nil {
		return fmt.Errorf("failed to open target file %s: %v", targetFile, err)
	}
	defer target.Close()

	if !utils.DestinationExists(destinationPath) {
		return fmt.Errorf("destination path (--dest) does not exist")
	}

	preEncoders, warnings := encoders.FromStrSlice(preEncoding)
	if warnings != "" {
		return fmt.Errorf("error determining pre-encoding: %v", warnings)
	}
	if len(preEncoders) > 5 {
		return fmt.Errorf("too many pre-encoders, limit 5")
	}

	return embedder.Process(&embedder.Config{
		Input:           input,
		SrcType:         srcType,
		SrcFilename:     getBaseFilename(target.Name()),
		Target:          target,
		DestinationPath: destinationPath,
		PreEncoding:     preEncoders,
	})
}

// getInputString determines if the input is either a string or another file
func getInputString(str string) (out, srcType string, err error) {
	info, err := os.Stat(inputStr)
	if err != nil {
		return str, "txt", nil
	}
	contents, err := os.ReadFile(str)
	if err != nil {
		return "", "", fmt.Errorf("error reading %s: %v", str, err)
	}
	srcType = filepath.Ext(info.Name())
	out = string(contents)
	return out, srcType, nil
}

func getBaseFilename(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}
