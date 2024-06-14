package path

import (
	"github.com/coderant163/docSyncKit/src/logger"
	"github.com/fsnotify/fsnotify"
)

// Monitor 监控目录变更
func Monitor(parentDir, repository string) error {
	watchDir, err := LocalPath(parentDir, repository)
	if err != nil {
		logger.Sugar().Errorf("LocalPath fail, err:%s", err.Error())
		return err
	}

	// Create a new file system watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Sugar().Errorf("fsnotify.NewWatcher fail, err:%s", err.Error())
	}
	defer watcher.Close()

	// Add the directory to be watched
	err = watcher.Add(watchDir)
	if err != nil {
		logger.Sugar().Errorf("watcher.Add fail, err:%s", err.Error())
	}

	// Listen for file system events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			// Handle different event types
			if event.Has(fsnotify.Create) {
				logger.Sugar().Infof("New file created: %s", event.Name)
			} else if event.Has(fsnotify.Write) {
				logger.Sugar().Infof("File modified: %s", event.Name)
			} else if event.Has(fsnotify.Remove) {
				logger.Sugar().Infof("File deleted: %s", event.Name)
			} else if event.Has(fsnotify.Rename) {
				logger.Sugar().Infof("File renamed: %s", event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return err
			}
			logger.Sugar().Errorf("watcher.Errors, err:%s", err.Error())
			return err
		}
	}

}
