package main

import (
	"github.com/coderant163/docSyncKit/src/cmd"
	"github.com/coderant163/docSyncKit/src/conf"
	"github.com/coderant163/docSyncKit/src/logger"
	"github.com/coderant163/docSyncKit/src/path"
)

func main() {
	logCfg := conf.Conf.Log
	logFile := path.GetFullLogFile(logCfg.FileName)
	logger.Init(logCfg.Level, logFile, logCfg.MaxSize, logCfg.MaxAge, logCfg.MaxBackups, logCfg.Compress)
	defer logger.Sync()

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
