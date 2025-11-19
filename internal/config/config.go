package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type KafkaTopics struct {
	MatchCreated string `yaml:"match_created"`
}
type KafkaConfig struct {
	Brokers []string    `yaml:"brokers"`
	GroupID string      `yaml:"groupId"`
	Topics  KafkaTopics `yaml:"topics"`
}
type Config struct {
	ServiceName string      `yaml:"serviceName"`
	HTTPPort    int         `yaml:"httpPort"`
	Kafka       KafkaConfig `yaml:"kafka"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
