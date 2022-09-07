package main

import (
	yaml "gopkg.in/yaml.v3"
	"os"
)

var configData *config

type config struct {
	Server struct {
		Bind string `yaml:"bind"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
}

func loadConfig(conf string) error {
	c := &config{}
	yamlFile, err := os.ReadFile(conf)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}
	configData = c
	return nil
}
