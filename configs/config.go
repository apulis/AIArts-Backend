package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type ProjectConfig struct {
	Port int
	Log  LogConfig
}

type LogConfig struct {
	WriteFile bool
	FileDir   string
	FileName  string
}

var Config ProjectConfig

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error read config file: %s \n", err))
	}

	viper.Unmarshal(&Config)
}
