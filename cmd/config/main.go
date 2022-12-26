package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Web3Config struct {
	Node    string `mapstructure:"node" json:"node" yaml:"node"`
	ChainId uint32 `mapstructure:"chainId" json:"chainId" yaml:"chainId"`
}

type Config struct {
	Web3 Web3Config `mapstructure:"web3" json:"web3" yaml:"web3"`
}

const (
	CONFIG_ENV  = "CONFIG_PATH"
	DEFAULT_ENV = "config.yaml"
)

// path > cmdline > dotenv > default
func GetConfigFile(path ...string) string {
	var file string
	if len(path) > 0 {
		file = path[0]
		fmt.Printf("using the func Viper(), config path is %s\n", file)
	} else {
		flag.StringVar(&file, "c", "", "choose config file.")
		flag.Parse()
		if file == "" {
			configEnv := os.Getenv(CONFIG_ENV)
			if configEnv == "" {
				file = DEFAULT_ENV
			} else {
				file = configEnv
				fmt.Printf("using env config, config path is %s\n", file)
			}
		} else {
			fmt.Printf("using the -c of the command line, config path is %s\n", file)
		}
	}
	return file
}

func NewViper(config *Config, path ...string) (*viper.Viper, error) {
	file := GetConfigFile(path...)
	v := viper.New()
	v.SetConfigFile(file)
	v.SetConfigType("yaml")
	fmt.Printf("file %s\n", file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(config); err != nil {
			fmt.Println(err)
		} else {
			showConfig(config)
		}
	})
	v.WatchConfig()

	if err := v.Unmarshal(config); err != nil {
		return nil, err
	}

	return v, nil
}

func showConfig(config *Config) {
	fmt.Printf("node: %s\n", config.Web3.Node)
	fmt.Printf("chainId: %d\n", config.Web3.ChainId)
}

func main() {
	config := &Config{}
	_, err := NewViper(config)
	if err != nil {
		panic(err.Error())
	}
	showConfig(config)
	select {}
}
