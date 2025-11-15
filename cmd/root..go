package cmd

import (
	"github.com/bshore/steggo/cmd/embed"
	"github.com/bshore/steggo/cmd/extract"

	"github.com/spf13/cobra"
)

const steggoLongHelpText = `
Try subcommands:

steggo embed --help
steggo extract --help
`

var rootCmd = &cobra.Command{
	Use:   "steggo",
	Short: "The base command of the LSB Steganography encode/decode CLI",
	Long:  steggoLongHelpText,
}

func Execute(args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func InitRoot() {
	embed.InitCmd()
	rootCmd.AddCommand(embed.Cmd)

	extract.InitCmd()
	rootCmd.AddCommand(extract.Cmd)
}
