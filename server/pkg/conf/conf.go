package conf

import (
	"supersign/pkg/tools"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type config struct {
	Server  server  `yaml:"SERVER"`
	Log     log     `yaml:"LOG"`
	Storage storage `yaml:"STORAGE"`
	Mysql   mysql   `yaml:"MYSQL"`
}

type server struct {
	URL          string `yaml:"URL"`
	MaxJob       int    `yaml:"MAXJOB"`
	RunMode      string `yaml:"RUNMODE"`
	ReadTimeout  int    `yaml:"READTIMEOUT"`
	WriteTimeout int    `yaml:"WRITETIMEOUT"`
	HttpPort     int    `yaml:"HTTPPORT"`
	TLS          bool   `yaml:"TLS"`
	Crt          string `yaml:"CRT"`
	Key          string `yaml:"KEY"`
}

type log struct {
	Level string `yaml:"LEVEL"`
}

type storage struct {
	EnableOSS          bool   `yaml:"ENABLEOSS"`
	BucketName         string `yaml:"BUCKETNAME"`
	OSSEndpoint        string `yaml:"OSSENDPOINT"`
	OSSAccessKeyId     string `yaml:"OSSACCESSKEYID"`
	OSSAccessKeySecret string `yaml:"OSSACCESSKEYSECRET"`
}

type mysql struct {
	Enable      bool   `yaml:"ENABLE"`
	Dsn         string `yaml:"DSN"`
	MaxIdle     int    `yaml:"MAXIDLE"`
	MaxOpen     int    `yaml:"MAXOPEN"`
	MaxLifetime int    `yaml:"MAXLIFETIME"`
}

type apple struct {
	AppleDeveloperPath string
	UploadFilePath     string
	TemporaryFilePath  string
}

var (
	Server  server
	Log     log
	Storage storage
	Mysql   mysql
	Apple   = apple{
		AppleDeveloperPath: "data/apple_developer/",
		UploadFilePath:     "data/upload_file_path/",
		TemporaryFilePath:  "data/temporary_file_path/",
	}
	Path string
)

// Setup 配置文件设置
func Setup() {
	if Path != "" {
		viper.SetConfigFile(Path)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("default")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := setConfig(); err != nil {
		panic(err)
	}
	mkdir([]string{Apple.AppleDeveloperPath, Apple.UploadFilePath, Apple.TemporaryFilePath})
}

func mkdir(paths []string) {
	for _, path := range paths {
		err := tools.MkdirAll(path)
		if err != nil {
			panic(err)
		}
	}
}

// Reset 配置文件重设
func Reset() error {
	return setConfig()
}

// OnChange 配置文件热加载回调
func OnChange(run func()) {
	viper.OnConfigChange(func(in fsnotify.Event) { run() })
	viper.WatchConfig()
}

// setConfig 构造配置文件到结构体对象上
func setConfig() error {
	var config config
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}
	Server = config.Server
	Log = config.Log
	Storage = config.Storage
	Mysql = config.Mysql
	return nil
}