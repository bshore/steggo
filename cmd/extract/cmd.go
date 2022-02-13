package extract

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "extract",
	Short: "",
	Long:  "",
	RunE:  extractCmdFn,
}

func extractCmdFn(command *cobra.Command, args []string) (err error) {
	return nil
}
