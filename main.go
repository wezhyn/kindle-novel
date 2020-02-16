package main

import (
	"flag"
	"kindle/config"
	M "kindle/data/email"
	"kindle/executor"
)

var (
	homeDir    string
	configName string
)

func init() {
	// -d 修改默认的配置文件地址
	flag.StringVar(&homeDir, "d", "", "set configuration directory")

	// -f 修改在配置目录下的文件名，例如：config.yaml
	flag.StringVar(&configName, "f", "", "specify configuration file")
	flag.Parse()
}

func main() {

	preInit()
	if cfg, configError := config.Parse(); configError != nil {
		panic(configError)
	} else {
		if refreshErr := config.RefreshConfig(*cfg); refreshErr != nil {
			panic(cfg)
		}
		// 初始化邮箱
		e := cfg.Email
		M.PostInit(e.Username, e.Password, e.Port, e.Host, e.Receiver, e.ErrReceiver)
		if exeErr := executor.Fun(*cfg); exeErr != nil {
			panic(exeErr)
		}
	}

}

func preInit() {
	if homeDir != "" {
		config.ConfigPath.ModifyHomeDir(homeDir)
	}
	if configName != "" {
		config.ConfigPath.ModifyConfigFile(configName)
	}
	defer func() {
		if err := config.Check(); err != nil {
			panic(err)
		}
	}()
}
