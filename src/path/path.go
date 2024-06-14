package path

import (
	"bufio"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coderant163/docSyncKit/src/logger"
	"golang.org/x/exp/maps"
)

const (
	MaxEncryptSize = 500
	pthSep         = string(os.PathSeparator)
)

func getCurrentDirectory() string {
	//返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Sugar().Fatalf("filepath.Abs failed![%v]\n", err)
	}

	//将\替换成/
	return strings.Replace(dir, "\\", pthSep, -1)
}

// GetConfDir 获取配置文件的目录
func GetConfDir() string {
	execDir := getCurrentDirectory()
	return fmt.Sprintf("%s%s..%sconf", execDir, pthSep, pthSep)
}

// getKeysDir 获取存放密钥的目录
func getKeysDir() string {
	execDir := getCurrentDirectory()
	return fmt.Sprintf("%s%s..%srsa_keys", execDir, pthSep, pthSep)
}

// GetKeyFilePath 获取存放密钥的目录
func GetKeyFilePath(keyFileName string) string {
	keysDir := getKeysDir()
	if _, err := CreateIfNotExists(keysDir); err != nil {
		logger.Sugar().Fatalf("CreateIfNotExists fail, err:[%s]", err.Error())
		return ""
	}
	return keysDir + pthSep + keyFileName
}

// GetFullLogFile 获取完整日志文件路径
func GetFullLogFile(fileName string) string {
	execDir := getCurrentDirectory()
	logDir := fmt.Sprintf("%s%s..%slogs", execDir, pthSep, pthSep)
	_, err := CreateIfNotExists(logDir)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s%s%s", logDir, pthSep, fileName)
}

func getLastSyncTimeFilePath() string {
	execDir := getCurrentDirectory()
	return fmt.Sprintf("%s%s..%sconf%ssync.txt", execDir, pthSep, pthSep, pthSep)
}

func LastSyncTime() (time.Time, error) {
	fileName := getLastSyncTimeFilePath()
	ok := Exists(fileName)
	if !ok {
		lastSync := time.Now()
		SetLastSyncTime(lastSync)
		return lastSync, nil
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		logger.Sugar().Errorf("os.ReadFile fail, err:%s", err.Error())
		return time.Time{}, err
	}
	preTime := strings.Split(string(content), "]")
	realTime := strings.Split(preTime[0], "[")

	lastSync, err := time.Parse(time.RFC3339, realTime[1])
	if err != nil {
		logger.Sugar().Errorf("time.Parse fail, err:%s", err.Error())
		return time.Time{}, err
	}
	return lastSync, nil
}

func SetLastSyncTime(lastSync time.Time) {
	fileName := getLastSyncTimeFilePath()
	timeStr := "[" + lastSync.Format(time.RFC3339) + "]"
	if err := os.WriteFile(fileName, []byte(timeStr), 0666); err != nil {
		logger.Sugar().Errorf("os.WriteFile fail, err:%s", err.Error())
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

// CreateIfNotExists 判断文件夹是否存在，如果不存在，创建该目录
func CreateIfNotExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, err
}

// GitPath 解析git地址，获取目录名称
// 示例： git@github.com:coderant163/docSyncKitTest.git
func gitPath(repository string) (string, error) {
	strs := strings.Split(repository, pthSep)
	if len(strs) == 0 {
		return "", errors.New("git地址不合法")
	}
	last := strs[len(strs)-1]
	strs = strings.Split(last, ".git")
	if len(strs) == 0 {
		return "", errors.New("git地址不合法")
	}
	return strs[0], nil
}

// SyncPath 主目录下用于执行git同步命令的隐藏目录
func SyncPath(parentDir string) string {
	return parentDir + pthSep + ".docSync"
}

// GitPath 本地用于执行git命令，与github进行同步数据的工作目录
func GitPath(parentDir, repository string) (string, error) {
	gitDir, err := gitPath(repository)
	if err != nil {
		return "", err
	}
	return SyncPath(parentDir) + pthSep + gitDir, nil
}

// LocalPath 本地文档路径
func LocalPath(parentDir, repository string) (string, error) {
	gitDir, err := gitPath(repository)
	if err != nil {
		return "", err
	}
	return parentDir + pthSep + gitDir, nil
}

// saveDirName 将加密后的目录名称保存到隐藏文件中
func saveDirName(dir, enc string) error {
	dstFile := dir + pthSep + ".dirInfo"
	// 打开目标文件，创建或者清空模式（实现覆盖写效果）
	dstFileHd, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		logger.Sugar().Errorf("os.OpenFile fail, err:%s", err.Error())
		return err
	}
	defer dstFileHd.Close()
	_, err = dstFileHd.WriteString(enc + "\n")
	if err != nil {
		logger.Sugar().Errorf("dstFileHd.WriteString [%s] fail, err:%s", enc, err.Error())
		return err
	}
	return nil
}

// loadDirName 读取加密后的目录名称
func loadDirName(dir string, keyStore KeyStore) (string, error) {
	dstFile := dir + pthSep + ".dirInfo"
	// 打开目标文件
	srcFileHd, err := os.OpenFile(dstFile, os.O_RDWR, 0666)
	if err != nil {
		logger.Sugar().Errorf("os.OpenFile fail, err:%s", err.Error())
		return "", err
	}
	defer srcFileHd.Close()
	r := bufio.NewReader(srcFileHd)
	buf, _, err := r.ReadLine()
	if err != nil {
		logger.Sugar().Errorf("srcFileHd.ReadLine fail, err:%s", err.Error())
		return "", err
	}
	newData, err := keyStore.Decrypt(string(buf))
	if err != nil {
		logger.Sugar().Errorf("keyStore.Decrypt [%s] fail, err:%s", string(buf), err.Error())
		return "", err
	}
	return string(newData), nil
}

// loadFileName 从加密文件中的第一行获取文件名称
func loadFileName(srcFile string, keyStore KeyStore) (string, error) {
	// 打开目标文件
	srcFileHd, err := os.OpenFile(srcFile, os.O_RDWR, 0666)
	if err != nil {
		logger.Sugar().Errorf("os.OpenFile fail, err:%s", err.Error())
		return "", err
	}
	defer srcFileHd.Close()
	r := bufio.NewReader(srcFileHd)
	buf, _, err := r.ReadLine()
	if err != nil {
		logger.Sugar().Errorf("srcFileHd.ReadLine fail, err:%s", err.Error())
		return "", err
	}
	newData, err := keyStore.Decrypt(string(buf))
	if err != nil {
		logger.Sugar().Errorf("keyStore.Decrypt [%s] fail, err:%s", string(buf), err.Error())
		return "", err
	}
	return string(newData), nil
}

type SyncFile struct {
	FileName    string //目标文件
	EncFileName string //加密后的文件名称
}

// ScanDir 遍历目录
//
// srcPath 原始目录
//
// dstPath 目标目录
//
// lastSync 上一次同步的时间
//
// isEnc 是否为加密流程。
// 若isEnc=true,为加密流程，需要将src修改时间晚于lastSync的文件加密，转存到dest目录
// 若isEnc=false,为解密流程，需要将src所有文件解密到dest目录
func ScanDir(srcPath, dstPath string, lastSync time.Time, isEnc bool, allowTypes []string,
	keyStore KeyStore) (fileMap map[string]SyncFile, err error) {
	fileMap = make(map[string]SyncFile)
	logger.Sugar().Debugf("srcPath:[%s],dstPath:[%s]", srcPath, dstPath)
	srcDir, err := os.ReadDir(srcPath)
	if err != nil {
		return nil, err
	}
	for _, fi := range srcDir {
		// 跳过隐藏文件或者隐藏目录
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		logger.Sugar().Debugf("scan [%s],isDir:%+v", fi.Name(), fi.IsDir())
		fileName := fi.Name()
		subSrcPath := srcPath + pthSep + fileName
		encFileName := ""  // 编码后的文件名称/路径名称
		showFileName := "" // 展示的名称
		if isEnc {
			encFileName, err = keyStore.Encrypt([]byte(fileName))
			if err != nil {
				logger.Sugar().Errorf("keyStore.Encrypt %s fail, err:%s", fileName, err.Error())
				return nil, err
			}
			showFileName = showName(fileName)
		}

		if fi.IsDir() { // 目录, 递归遍历
			if !isEnc {
				showFileName, err = loadDirName(subSrcPath, keyStore)
				if err != nil {
					logger.Sugar().Errorf("loadDirName %s fail, err:%s", subSrcPath, err.Error())
					return nil, err
				}
			}

			subDstPath := dstPath + pthSep + showFileName

			logger.Sugar().Infof("scan dir:%s", fileName)
			_, err = CreateIfNotExists(subDstPath)
			if err != nil {
				logger.Sugar().Errorf("CreateIfNotExists %s fail, err:%s", subDstPath, err.Error())
				return nil, err
			}
			if isEnc {
				//将编码后的文件名称写入到该目录下的隐藏文件中
				err = saveDirName(subDstPath, encFileName)
				if err != nil {
					logger.Sugar().Errorf("saveDirName %s fail, err:%s", subSrcPath, err.Error())
					return nil, err
				}
			}

			subFileMap, err := ScanDir(subSrcPath, subDstPath, lastSync, isEnc, allowTypes, keyStore)
			if err != nil {
				logger.Sugar().Errorf("ScanDir %s fail, err:%s", subSrcPath, err.Error())
				return nil, err
			} else {
				//合并子目录的文件
				maps.Copy(fileMap, subFileMap)
			}
		} else {
			ok := checkFileType(fi.Name(), isEnc, allowTypes)
			if ok {
				fileInfo, err := os.Stat(subSrcPath)
				if err != nil {
					logger.Sugar().Errorf("os.Stat %s fail, err:%s", subSrcPath, err.Error())
					return nil, err
				} else {
					if !isEnc {
						// 从文件中加载真实文件名称
						showFileName, err = loadFileName(subSrcPath, keyStore)
						if err != nil {
							logger.Sugar().Errorf("loadFileName %s fail, err:%s", subSrcPath, err.Error())
							return nil, err
						}
					}
					subDstPath := dstPath + pthSep + showFileName
					logger.Sugar().Debugf("file [%s] modTime is [%s], lastSync is [%s] ,After:%+v",
						fileInfo.Name(),
						fileInfo.ModTime().String(), lastSync.String(), fileInfo.ModTime().After(lastSync))
					if !isEnc || (isEnc && fileInfo.ModTime().After(lastSync)) {
						fileMap[subSrcPath] = SyncFile{FileName: subDstPath, EncFileName: encFileName}
					}
				}
			} else {
				logger.Sugar().Infof("check file [%s] type fail", fi.Name())
			}
		}
	}
	return fileMap, nil
}

func checkFileType(fileName string, isEnc bool, allowTypes []string) bool {
	if !isEnc {
		return true
	}

	ok := false
	for i := 0; i < len(allowTypes); i++ {
		if allowTypes[i] == "*" {
			return true
		}
		ok = strings.HasSuffix(fileName, allowTypes[i])
		if ok {
			break
		}
	}
	return ok
}

type KeyStore interface {
	Encrypt(plainText []byte) (string, error)
	Decrypt(text string) ([]byte, error)
}

func TransferFile(srcFile string, dstFileInfo SyncFile, isEnc bool, keyStore KeyStore) error {
	// 打开原文件，只读模式
	srcFileHd, err := os.OpenFile(srcFile, os.O_RDWR, 0666)
	if err != nil {
		logger.Sugar().Errorf("os.OpenFile fail, err:%s", err.Error())
		return err
	}
	defer srcFileHd.Close()

	_, err = srcFileHd.Stat()
	if err != nil {
		logger.Sugar().Errorf("srcFileHd.Stat fail, err:%s", err.Error())
		return err
	}

	// 打开目标文件，创建或者清空模式（实现覆盖写效果）
	dstFileHd, err := os.OpenFile(dstFileInfo.FileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		logger.Sugar().Errorf("os.OpenFile fail, err:%s", err.Error())
		return err
	}
	defer dstFileHd.Close()

	if isEnc {
		//写入文件名称
		_, err = dstFileHd.WriteString(dstFileInfo.EncFileName + "\n")
		if err != nil {
			logger.Sugar().Errorf("dstFileHd.WriteString [%s] fail, err:%s", dstFileInfo.EncFileName, err.Error())
			return err
		}
		return handleEncryptFile(srcFileHd, dstFileHd, keyStore)
	}
	return handleDecryptFile(srcFileHd, dstFileHd, keyStore)
}
func handleEncryptFile(srcFileHd, dstFileHd *os.File, keyStore KeyStore) error {
	for {
		buf := make([]byte, MaxEncryptSize)
		n, err := srcFileHd.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				logger.Sugar().Errorf("srcFileHd.Read fail, err:%s", err.Error())
				return err
			}
		}

		newData, err := keyStore.Encrypt(buf[:n])
		if err != nil {
			logger.Sugar().Errorf("keyStore.Encrypt [%s] fail, err:%s", string(buf), err.Error())
			return err
		}
		_, err = dstFileHd.WriteString(newData + "\n")
		if err != nil {
			logger.Sugar().Errorf("dstFileHd.WriteString [%s] fail, err:%s", newData, err.Error())
			return err
		}
	}
	return nil
}

func handleDecryptFile(srcFileHd, dstFileHd *os.File, keyStore KeyStore) error {
	r := bufio.NewReader(srcFileHd)
	isFirstLine := true
	for {
		// ReadLine is a low-level line-reading primitive.
		// Most callers should use ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
		buf, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Sugar().Errorf("srcFileHd.ReadLine fail, err:%s", err.Error())
			return err
		}
		// 过滤第一行数据，该内容为文件名称
		if isFirstLine {
			isFirstLine = false
			continue
		}
		newData, err := keyStore.Decrypt(string(buf))
		if err != nil {
			logger.Sugar().Errorf("keyStore.Decrypt [%s] fail, err:%s", string(buf), err.Error())
			return err
		}
		_, err = dstFileHd.Write(newData)
		if err != nil {
			logger.Sugar().Errorf("dstFileHd.Write [%s] fail, err:%s", newData, err.Error())
			return err
		}
	}
	return nil
}

func showName(s string) string {
	sum := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", sum)
}
