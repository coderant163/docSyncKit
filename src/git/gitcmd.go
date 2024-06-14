package git

import (
	"fmt"
	"os/exec"

	"github.com/coderant163/docSyncKit/src/logger"
	"github.com/coderant163/docSyncKit/src/path"
)

type Client struct {
	SyncPath   string //执行git clone命令的目录
	GitDir     string //clone后的仓库路径
	Repository string //仓库地址ssh
	Branch     string // 分支名称
	Name       string // 用户名
	Email      string // 邮箱名称
}

func NewClient(parentDir, repository, branch, name, email string) (*Client, error) {
	gitDir, err := path.GitPath(parentDir, repository)
	if err != nil {
		logger.Sugar().Errorf("path.GitPath fail, err:%s", err.Error())
		return nil, err
	}
	syncPath := path.SyncPath(parentDir)
	// 如果syncPath目录不存在，创建
	_, err = path.CreateIfNotExists(syncPath)
	if err != nil {
		logger.Sugar().Errorf("path.CreateIfNotExists fail, err:%s", err.Error())
		return nil, err
	}

	return &Client{
		SyncPath:   syncPath,
		GitDir:     gitDir,
		Repository: repository,
		Branch:     branch,
		Name:       name,
		Email:      email,
	}, nil
}

// Clone 从远端克隆项目
func (c *Client) Clone() {
	if path.Exists(c.GitDir) {
		logger.Sugar().Infof("已经存在clone目标路径%s，跳过该步骤", c.GitDir)
		return
	}
	stdout, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && git clone %s", c.SyncPath,
		c.Repository)).Output()
	if err != nil {
		logger.Sugar().Errorf("exec.Command fail, err:%s", err.Error())
	}
	logger.Sugar().Infof("git clone output:[%s]", string(stdout))
}

// InitGitConfig 初始化git项目的用户信息
func (c *Client) InitGitConfig() {
	stdout, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && "+
		"git config user.name %s && "+
		"git config user.email %s",
		c.GitDir,
		c.Name,
		c.Email)).Output()
	if err != nil {
		logger.Sugar().Errorf("exec.Command fail, err:%s", err.Error())
	}
	logger.Sugar().Infof("git config output:[%s]", string(stdout))
}

// Checkout 检出分支，暂时不用
func (c *Client) Checkout() {
	stdout, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && git checkout  %s", c.GitDir,
		c.Branch)).Output()
	if err != nil {
		logger.Sugar().Errorf("exec.Command fail, err:%s", err.Error())
	}
	logger.Sugar().Infof("git checkout output:[%s]", string(stdout))
}

// CommitAll 提交所有改动，并推送到远端仓库
func (c *Client) CommitAll() {
	stdout, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && git add * && "+
		"git commit -m 'auto commit' && "+
		"git push origin main",
		c.GitDir)).Output()
	if err != nil {
		logger.Sugar().Errorf("exec.Command fail, err:%s", err.Error())
	}
	logger.Sugar().Infof("git commit output:[%s]", string(stdout))
}
