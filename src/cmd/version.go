package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	CommitID  = "2"
	Branch    = "2"
	BuildTime = "2"
)

func init() {
	rootCmd.AddCommand(NewVersionCommand())
}

// NewVersionCommand 监控本地目录中的文件变化
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: `version info`,
		Long:  `version info`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("git branch:%s\n", Branch)
			fmt.Printf("git commitID:%s\n", CommitID)
			fmt.Printf("build time:%s\n", BuildTime)
			return nil
		},
	}
	return cmd
}
