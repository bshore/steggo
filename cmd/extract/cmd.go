package extract

import (
	"fmt"
	"os"

	"github.com/bshore/steggo/pkg/extractor"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "extract",
	Short: "Extracts a message from --target {file} outputting to --dest {path}",
	RunE:  extractCmdFn,
}

var (
	targetFile      string
	destinationPath string
)

func InitCmd() {
	Cmd.PersistentFlags().StringVarP(&targetFile, "target", "t", "", "The path to the image file being targeted for extraction")
	Cmd.PersistentFlags().StringVarP(&destinationPath, "dest", "d", "", "The destination path to output the extracted message (message.txt)")
}

func extractCmdFn(command *cobra.Command, args []string) (err error) {
	target, err := os.Open(targetFile)
	if err != nil {
		return fmt.Errorf("failed to open target file %s: %v", targetFile, err)
	}
	defer target.Close()
	return extractor.Process(&extractor.Config{
		Target:          target,
		DestinationPath: destinationPath,
	})
}
