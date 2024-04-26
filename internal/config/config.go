package config

import (
	"flag"
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
	// TODO: implement config parsing from file
	return nil
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
