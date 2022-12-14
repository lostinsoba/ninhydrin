package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const configPath = "/etc/ninhydrin/ninhydrin.yml"

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	return loadFromBytes(data)
}

func loadFromBytes(data []byte) (cfg *Config, err error) {
	err = yaml.Unmarshal(data, &cfg)
	return
}

type Config struct {
	API        API        `yaml:"api"`
	Scheduler  Scheduler  `yaml:"scheduler"`
	Storage    Storage    `yaml:"storage"`
	Monitoring Monitoring `yaml:"monitoring"`
}

type API struct {
	Listen string `yaml:"listen"`
}

type Scheduler struct {
	Interval time.Duration `yaml:"interval"`
}

type Storage struct {
	Kind     string            `yaml:"kind"`
	Settings map[string]string `yaml:"settings"`
}

type Monitoring struct {
	Logger   Logger   `yaml:"logger"`
	Exporter Exporter `yaml:"exporter"`
}

type Logger struct {
	Kind     string            `yaml:"kind"`
	Settings map[string]string `yaml:"settings"`
}

type Exporter struct {
	Kind     string            `yaml:"kind"`
	Settings map[string]string `yaml:"settings"`
}
