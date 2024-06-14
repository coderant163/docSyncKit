package cmd

import (
	"github.com/coderant163/docSyncKit/src/conf"
	"github.com/coderant163/docSyncKit/src/logger"
	"github.com/coderant163/docSyncKit/src/rsa"
	"github.com/spf13/cobra"
)

var (
	data = ""
)

func init() {
	rsaCmd := &cobra.Command{
		Use:   "rsa",
		Short: `rsa create/encrypt/decrypt`,
		Long:  `rsa create/encrypt/decrypt`,
	}
	rsaCmd.AddCommand(NewRSACreateCommand())
	encryptCmd := NewRSAEncryptCommand()
	encryptCmd.Flags().StringVarP(&data, "data", "d", "", "data")
	rsaCmd.AddCommand(encryptCmd)
	decryptCmd := NewRSADecryptCommand()
	decryptCmd.Flags().StringVarP(&data, "data", "d", "", "data")
	rsaCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(rsaCmd)
}

func NewRSACreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: `rsa create key`,
		Long:  `rsa create key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Sugar().Infof("run rsa create cmd")

			err := rsa.GenerateRSAKey()
			if err != nil {
				logger.Sugar().Fatalf("rsa.GenerateRSAKey fail, err:[%s]", err.Error())
			}
			logger.Sugar().Infof("create rsa keys success")
			return nil
		},
	}
	return cmd
}

func NewRSAEncryptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: `rsa encrypt msg`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Sugar().Infof("run rsa encrypt cmd")

			keyStore := rsa.NewRSA(conf.Conf.Base.PrivateKeyFile, conf.Conf.Base.PublicKeyFile)
			text, err := keyStore.Encrypt([]byte(data))
			if err != nil {
				logger.Sugar().Fatalf("keyStore.Encrypt fail, err:[%s]", err.Error())
			}
			logger.Sugar().Infof("keyStore.Encrypt success, text:[%s]", text)
			return nil
		},
	}
	return cmd
}

func NewRSADecryptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrypt",
		Short: `rsa decrypt msg`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Sugar().Infof("run rsa decrypt cmd")

			keyStore := rsa.NewRSA(conf.Conf.Base.PrivateKeyFile, conf.Conf.Base.PublicKeyFile)
			text, err := keyStore.Decrypt(data)
			if err != nil {
				logger.Sugar().Fatalf("keyStore.Decrypt fail, err:[%s]", err.Error())
			}
			logger.Sugar().Infof("keyStore.Decrypt success, text:[%s]", string(text))
			return nil
		},
	}
	return cmd
}
