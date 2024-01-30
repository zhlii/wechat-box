package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Common struct {
		LogLevel string `yaml:"log_level"`
		IsProd   bool   `yaml:"is_prod"`
	}

	Httpd struct {
		Addr   string
		Enable bool
	}

	Rpc struct {
		Host   string
		Port   int
		Enable bool
	}
}

var Data *Config

func Init(path string) {
	fmt.Printf("read config from: %s\n", path)

	conf := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("read config error :%v\n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		fmt.Printf("yaml decode error :%v", err)
		os.Exit(1)
	}

	Data = conf
}
