package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// 检查 config.ConfigPath 包变量对应的路径是否存在(不尝试创建)
func Check() error {
	if _, err := os.Stat(ConfigPath.homeDir); os.IsNotExist(err) {
		return fmt.Errorf("无配置文件目录：%s : %s ", ConfigPath.homeDir, err.Error())
	}
	if _, err := os.Stat(ConfigPath.ConfigFile()); err != nil {
		return fmt.Errorf("配置文件出错：%s : %s", ConfigPath.configFile, err.Error())
	}
	if ext := filepath.Ext(ConfigPath.configFile); ext != ".yml" {
		return fmt.Errorf("配置文件格式出错，应该使用：.yml后缀,不是：%s", ext)
	}
	return nil
}
