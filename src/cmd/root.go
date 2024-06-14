package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "docSyncKit",
	Short: "docSyncKit",
	Long:  `doc sync kit`,
}

func init() {

}

func Execute() error {
	return rootCmd.Execute()
}
