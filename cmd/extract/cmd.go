package extract

import (
	"fmt"
	"os"

	"github.com/bshore/steggo/pkg/extractor"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "extract",
	Short: "",
	Long:  "",
	RunE:  extractCmdFn,
}

var (
	targetFile      string
	destinationPath string
)

func InitCmd() {
	Cmd.PersistentFlags().StringVar(&targetFile, "target", "", "The path to the image file being targeted for extraction")
	Cmd.PersistentFlags().StringVar(&destinationPath, "dest", "", "The destination path to output the extracted message (message.txt)")
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
