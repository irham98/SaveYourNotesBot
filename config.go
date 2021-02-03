package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Credentials string `yaml:"credentials"`
	Port        int32  `yaml:"port"`
	Telegram    struct {
		WebhookUrl string `yaml:"webhook_url"`
		Token      string `yaml:"token"`
	} `yaml:"telegram"`
}

func loadConfig(configFile string) (config *Config, err error) {
	configContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config = &Config{}
	err = yaml.Unmarshal(configContent, config)
	if err != nil {
		return nil, err
	}

	return
}
