package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type APIconfig struct {
	BaseUrl string `yaml:"BASE_URL"`
	APIkey  string `yaml:"API_KEY"`
    Timeout time.Duration
}

func NewConfig(timeout time.Duration) (*APIconfig, error) {
	var config *APIconfig = &APIconfig{}
	configPath, err := GetconfigPath()
	defer func() {
		if err != nil {
			slog.Error("cannot set up config:invalid path")
		}
	}()
	f, err := os.Open(configPath)

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	if err != nil {
		return nil, err
	}
    if timeout != 0 {
        config.Timeout = timeout
    }

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(config); err != nil {
	}
	return config, nil

}

func GetconfigPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, "config.yaml"), nil
}
