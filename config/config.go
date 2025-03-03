package config

import (
	"gopkg.in/yaml.v2"
	"log/slog"
	"os"
	"path/filepath"
)

type APIconfig struct {
	BaseUrl string `yaml:"BASE_URL"`
	APIkey  string `yaml:"API_KEY"`
}

func NewConfig() (*APIconfig, error) {
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
