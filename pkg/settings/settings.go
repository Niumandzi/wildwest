package settings

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	API struct {
		Port int `yaml:"port"`
	} `yaml:"api"`
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
}

func (c *Config) ReadConfig() error {
	yamlFile, err := os.ReadFile("configs/dev.yaml")
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		return err
	}

	return nil
}
