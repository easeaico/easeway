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

type Email struct {
	Subject      string `yaml:"subject"`
	From         string `yaml:"from"`
	Pass         string `yaml:"pass"`
	ProviderHost string `yaml:"provider_host"`
	ProviderPort int    `yaml:"provider_port"`
}
type Config struct {
	Email        Email  `yaml:"email"`
	Gemini       Gemini `yaml:"gemini"`
	OpenAI       OpenAI `yaml:"openai"`
	DBConnection string `yaml:"db_connection"`
	Server       Server `yaml:"server"`
	Site         Site   `yaml:"site"`
}

type Site struct {
	Scheme string `yaml:"scheme"`
	Host   string `yaml:"host"`
}

var Conf Config

func NewConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading YAML file: %s", err)
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML data: %s", err)
	}

	return &Conf
}
