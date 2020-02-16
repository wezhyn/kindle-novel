package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

// 包级别变量，在 init() 中初始化
var ConfigPath *Path

var once sync.Once

var globalConfig *Config

//默认~/.config/kindle/config.yml
func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir, _ = os.Getwd()
	}
	homeDir = path.Join(homeDir, ".config", "kindle")
	configFile := "config.yml"
	ConfigPath = &Path{homeDir: homeDir, configFile: configFile}
}

func NewConfig() *Config {
	once.Do(func() {
		globalConfig = new(Config)
	})
	return globalConfig
}

// 从配置文件中解析Config
func Parse() (*Config, error) {
	c := new(Config)
	if err := Check(); err != nil {
		return nil, err
	}
	buf, err := readConfig()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func RefreshConfig(config Config) error {
	originConfig := NewConfig()
	originConfig.UpdateUserAgent(config.UserAgent)
	originConfig.UpdateNovels(config.Novels)
	return nil
}

func readConfig() ([]byte, error) {
	file, err := ioutil.ReadFile(ConfigPath.ConfigFile())
	if err == nil && len(file) != 0 {
		return file, nil
	}
	return nil, fmt.Errorf("配置文件无数据：%s", ConfigPath.ConfigFile())
}
