package config

import (
	"log"
	"os"

	"go-seven/pkg/models"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/conf"
)

var (
	vp     *viper.Viper
	config models.Config
)

const (
	CONFIG_ENV  = "CONFIG_PATH"
	DEFAULT_ENV = "config.yaml"
)

// path > cmdline > dotenv > default
func GetConfigFile(path ...string) string {
	var file string
	if len(path) > 0 {
		file = path[0]
	} else {
		configEnv := os.Getenv(CONFIG_ENV)
		if configEnv == "" {
			file = DEFAULT_ENV
		} else {
			file = configEnv
		}
	}
	return file
}

func NewViper(path ...string) error {
	file := GetConfigFile(path...)
	v := viper.New()
	v.SetConfigFile(file)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&config); err != nil {
			log.Fatalln(err.Error())
		}
	})
	v.WatchConfig()

	if err := v.Unmarshal(&config); err != nil {
		return err
	}

	vp = v
	return nil
}

func GetViper() *viper.Viper {
	return vp
}

func GetConfig() *models.Config {
	return &config
}

func New(path ...string) error {
	godotenv.Load()
	file := GetConfigFile(path...)
	return conf.Load(file, &config, conf.UseEnv())
}
