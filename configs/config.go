package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type ProjectConfig struct {
	Port int
	Log  LogConfig
	Db   DbConfig
	File FileConfig
<<<<<<< HEAD
	Auth AuthConfig
=======
	DltsUrl string
>>>>>>> 4d357c1d0358af2890ab724dbd475f142614556a
}

type LogConfig struct {
	WriteFile bool
	FileDir   string
	FileName  string
}

type DbConfig struct {
	Username     string
	Password     string
	Host         string
	Port         int
	Database     string
	MaxOpenConns int
	MaxIdleConns int
}

type FileConfig struct {
	DatasetDir string
	ModelDir   string
	SizeLimit  int
}

type AuthConfig struct {
	Url string
	Key string
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
