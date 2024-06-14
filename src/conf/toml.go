package conf

import (
	"github.com/coderant163/docSyncKit/src/path"
	"github.com/spf13/viper"
)

var Conf *Config

func init() {
	c, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	Conf = c
}

// Config 配置
type Config struct {
	Base   BaseCfg `toml:"Base"`
	Github Github  `toml:"Github"`
	Log    Log     `toml:"Log"`
}

type BaseCfg struct {
	WorkDir        string   `toml:"WorkDir"`        // 工作目录
	AllowFileType  []string `toml:"AllowFileType"`  // 允许同步到云端的文件类型
	PrivateKeyFile string   `toml:"PrivateKeyFile"` // 私钥文件名，从云端同步到本地时需要该参数
	PublicKeyFile  string   `toml:"PublicKeyFile"`  // 公钥文件名，推送到云端需要该参数
}

type Github struct {
	Repository string `toml:"Repository"` // 仓库地址,ssh方式
	Branch     string `toml:"Branch"`     // 默认分支，本参数暂时没用到
	Name       string `toml:"Name"`       // git用户名
	Email      string `toml:"Email"`      // git邮箱名
}

type Log struct {
	Level      string `toml:"Level"`      // 日志级别
	FileName   string `toml:"FileName"`   // 日志名称
	MaxSize    int    `toml:"MaxSize"`    // 日志大小限制，单位MB
	MaxAge     int    `toml:"MaxAge"`     // 历史日志文件保留天数
	MaxBackups int    `toml:"MaxBackups"` // 最大保留历史日志数量
	Compress   bool   `toml:"Compress"`   // 历史日志文件压缩标识
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	viper.SetConfigName("conf")            // name of config file (without extension)
	viper.SetConfigType("toml")            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path.GetConfDir()) // path to look for the config file in
	//viper.AddConfigPath("../conf") // call multiple times to add many search paths
	//viper.AddConfigPath(".")       // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return nil, err
	}
	var conf Config

	err = viper.Unmarshal(&conf)
	return &conf, err
}
