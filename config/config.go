package config

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
)

var configFileName string

type Config struct {
	LogfileName         string `json:"logfile_name"`
	LogLevel            string `json:"logfile_level"`
	MaximumConcurrently int    `json:"max_requests_concurrently"`
	MaximumRPM          int    `json:"max_rpm"`
	MaxURISize          int    `json:"max_uri_size"`
}

var defaultCfg = &Config{
	LogLevel:            "trace",
	MaximumConcurrently: 10,
	MaximumRPM:          50,
}

func ReadConfig() (*Config, error) {
	if configFileName == "" {
		return defaultCfg, nil
	}

	buf, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, errors.New("config: " + err.Error())
	}

	var cfg Config
	if err := json.Unmarshal(buf, &cfg); err != nil {
		return nil, errors.New("config: " + err.Error())
	}
	return &cfg, nil
}

func init() {
	flag.StringVar(&configFileName, "c", "", "config filename")
}
