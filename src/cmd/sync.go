package cmd

import (
	"time"

	"github.com/coderant163/docSyncKit/src/conf"
	"github.com/coderant163/docSyncKit/src/git"
	"github.com/coderant163/docSyncKit/src/logger"
	"github.com/coderant163/docSyncKit/src/path"
	"github.com/coderant163/docSyncKit/src/rsa"
	"github.com/spf13/cobra"
)

func init() {
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: `sync to remote/local`,
		Long:  `sync to remote/local`,
	}

	localCmd := NewSyncToLocalCommand()
	remoteCmd := NewSyncToRemoteCommand()
	//rsaCmd.Flags().StringVarP(&srcData, "data", "", "", "data")
	syncCmd.AddCommand(localCmd)
	syncCmd.AddCommand(remoteCmd)
	rootCmd.AddCommand(syncCmd)
}

func NewSyncToLocalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: `sync from remote to local`,
		Long:  `sync from remote to local`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Sugar().Infof("run sync from remote to local cmd")
			logger.Sugar().Infof("dir:%s,github:%s,branch:%s,name:%s,email:%s", conf.Conf.Base.WorkDir,
				conf.Conf.Github.Repository, conf.Conf.Github.Branch, conf.Conf.Github.Name, conf.Conf.Github.Email)

			gitClient, err := git.NewClient(conf.Conf.Base.WorkDir, conf.Conf.Github.Repository,
				conf.Conf.Github.Branch,
				conf.Conf.Github.Name, conf.Conf.Github.Email)
			if err != nil {
				logger.Sugar().Errorf("git.NewClient fail, err:%s", err.Error())
				return err
			}
			gitClient.Clone()
			logger.Sugar().Infof("git clone ok")
			gitClient.InitGitConfig()
			logger.Sugar().Infof("git init ok")
			isEnc := false
			srcPath, _ := path.GitPath(conf.Conf.Base.WorkDir, conf.Conf.Github.Repository)
			dstPath, _ := path.LocalPath(conf.Conf.Base.WorkDir, conf.Conf.Github.Repository)
			lastSyncTime := time.Now()

			_, err = path.CreateIfNotExists(dstPath)
			if err != nil {
				logger.Sugar().Errorf("CreateIfNotExists %s fail, err:%s", dstPath, err.Error())
				return err
			}
			keyStore := rsa.NewRSA(conf.Conf.Base.PrivateKeyFile, conf.Conf.Base.PublicKeyFile)
			fileMap, err := path.ScanDir(srcPath, dstPath, lastSyncTime, isEnc, conf.Conf.Base.AllowFileType, keyStore)
			if err != nil {
				logger.Sugar().Errorf("path.ScanDir fail, err:%s", err.Error())
				return err
			}
			for k, v := range fileMap {
				logger.Sugar().Infof("begin sync file from [%s] to [%s]", k, v)

				err = path.TransferFile(k, v, isEnc, keyStore)
				if err != nil {
					logger.Sugar().Errorf("path.TransferFile fail, err:%s", err.Error())
					return err
				}
				logger.Sugar().Infof("sync file from [%s] to [%s] success", k, v)
			}
			path.SetLastSyncTime(lastSyncTime)
			logger.Sugar().Infof("run sync from remote to local cmd success")
			return nil
		},
	}
	return cmd
}

func NewSyncToRemoteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote",
		Short: `sync from local to remote`,
		Long:  `sync from local to remote`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Sugar().Infof("run sync from local to remote cmd")
			logger.Sugar().Infof("dir:%s,github:%s,branch:%s,name:%s,email:%s", conf.Conf.Base.WorkDir,
				conf.Conf.Github.Repository, conf.Conf.Github.Branch, conf.Conf.Github.Name, conf.Conf.Github.Email)

			isEnc := true
			srcPath, _ := path.LocalPath(conf.Conf.Base.WorkDir, conf.Conf.Github.Repository)
			dstPath, _ := path.GitPath(conf.Conf.Base.WorkDir, conf.Conf.Github.Repository)
			timeNow := time.Now()
			lastSyncTime, err := path.LastSyncTime()
			if err != nil {
				logger.Sugar().Errorf("path.LastSyncTime fail, err:%s", err.Error())
				return err
			}
			logger.Sugar().Infof("lastSyncTime is %s", lastSyncTime.String())
			keyStore := rsa.NewRSA(conf.Conf.Base.PrivateKeyFile, conf.Conf.Base.PublicKeyFile)
			fileMap, err := path.ScanDir(srcPath, dstPath, lastSyncTime, isEnc, conf.Conf.Base.AllowFileType, keyStore)
			if err != nil {
				logger.Sugar().Errorf("path.ScanDir fail, err:%s", err.Error())
				return err
			}
			for k, v := range fileMap {
				logger.Sugar().Infof("begin sync file from [%s] to [%s]", k, v)

				err = path.TransferFile(k, v, isEnc, keyStore)
				if err != nil {
					logger.Sugar().Errorf("path.TransferFile fail, err:%s", err.Error())
					return err
				}
				logger.Sugar().Infof("sync file from [%s] to [%s] success", k, v)
			}
			if len(fileMap) > 0 {
				gitClient, err := git.NewClient(conf.Conf.Base.WorkDir, conf.Conf.Github.Repository,
					conf.Conf.Github.Branch,
					conf.Conf.Github.Name, conf.Conf.Github.Email)
				if err != nil {
					logger.Sugar().Errorf("git.NewClient fail, err:%s", err.Error())
					return err
				}
				gitClient.CommitAll()
			}

			path.SetLastSyncTime(timeNow)
			logger.Sugar().Infof("run sync from local to remote cmd success")
			return nil
		},
	}
	return cmd
}
