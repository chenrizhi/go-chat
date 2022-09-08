package main

import (
	yaml "gopkg.in/yaml.v3"
	"os"
	"time"
)

var configData *config

type config struct {
	Server struct {
		Bind         string `yaml:"bind"`
		Port         int    `yaml:"port"`
		AliveTimeout string `yaml:"aliveTimeout"`
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
	if _, err := time.ParseDuration(c.Server.AliveTimeout); err != nil {
		return err
	}
	configData = c
	return nil
}
