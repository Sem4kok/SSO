package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func Load() *Config {
	// get path to config.yaml
	// if it hasn't specified then panic
	path := getConfigPath()
	if path == "" {
		panic("invalid path to config file")
	}

	// check for file-existing (config.yaml)
	// if it isn't then panic
	if _, err := os.Stat(path); err != nil {
		panic("config file with specified path: " + path + " does not exist")
	}

	// decode config file into struct Config
	var cfg Config
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("failed ro read config: " + err.Error())
	}

	return &cfg
}

// getConfigPath return path to config
// Priority of get path: flag > env > default
// For local Load's is good to run app with:
//
//	flag: $ sso --config="path/to/config.yaml"
//	env:  $ CONFIG_PATH="path/to/config.yaml" sso
func getConfigPath() string {
	var path string

	// try to parse from flag
	flag.StringVar(&path, "config", "", "path to config")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
