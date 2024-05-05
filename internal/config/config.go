package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type OpenAI struct {
	ApiKey string `yaml:"api_key"`
}

type Gemini struct {
	ApiKey string `yaml:"api_key"`
}

type Server struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type Config struct {
	Gemini       Gemini `yaml:"gemini"`
	OpenAI       OpenAI `yaml:"openai"`
	DBConnection string `yaml:"db_connection"`
	Server       Server `yaml:"server"`
}

func NewConfig(path string) *Config {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading YAML file: %s", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML data: %s", err)
	}

	return &config
}
