package cmd

import (
	"github.com/coderant163/docSyncKit/src/conf"
	"github.com/coderant163/docSyncKit/src/logger"
	"github.com/coderant163/docSyncKit/src/path"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewMonitorCommand())
}

// NewMonitorCommand 监控本地目录中的文件变化
func NewMonitorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: `monitor local WorkDir, and sync changes from  local to remote`,
		Long:  `monitor local WorkDir, and sync changes from  local to remote`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Sugar().Infof("run monitor cmd")
			err := path.Monitor(conf.Conf.Base.WorkDir, conf.Conf.Github.Repository)
			if err != nil {
				logger.Sugar().Errorf("path.Monitor fail, err:%s", err.Error())
				return err
			}
			logger.Sugar().Infof("run monitor cmd success")
			return nil
		},
	}
	return cmd
}
