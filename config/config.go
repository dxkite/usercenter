package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	EnableUser bool   `yaml:"enable_user"`
	DataPath   string `yaml:"data_path"`
	Prefix     string `yaml:"api_prefix"`
}

func (cfg *Config) GetPrefix() string {
	if len(cfg.Prefix) == 0 {
		return "/api/user"
	}
	return cfg.Prefix
}

func (cfg *Config) LoadFrom(in []byte) error {
	if err := yaml.Unmarshal(in, cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) LoadFromFile(p string) error {
	if b, err := ioutil.ReadFile(p); err == nil {
		return cfg.LoadFrom(b)
	} else {
		return err
	}
}
