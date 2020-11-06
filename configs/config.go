package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type ProjectConfig struct {
	Port               int
	Log                LogConfig
	Db                 DbConfig
	File               FileConfig
	Auth               AuthConfig
	DltsUrl            string
	Anno               AnnotationConfig
	Infer              InferenceConfig
	Token              string
	Image              map[string][]string
	InteractiveModeJob bool
	PrivateRegistry    string
	Imagesave          ImageSaveConfig
	TrackingUrl        string
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
	DatasetDir         string
	ModelDir           string
	SizeLimit          int
	CleanBeforeSeconds int64
	CleanEverySeconds  int64
}

type AnnotationConfig struct {
	BackendUrl string
}
type InferenceConfig struct {
	BackendUrl string
}

type AuthConfig struct {
	Url                string
	Key                string
	SamlIdpMetadataURL string
	SamlRootUrl        string
	SamlPrivateKey     string
	SamlCertificate    string
}

type ImageSaveConfig struct {
	K8sconf string
	Sshkey  string
	Sshuser string
	Sshport string
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
	if Config.File.CleanBeforeSeconds <= 0 {
		Config.File.CleanBeforeSeconds = 86400
	}
	if Config.File.CleanEverySeconds <= 0 {
		Config.File.CleanEverySeconds = 600
	}
}
