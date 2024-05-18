package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	ErrFailedToReadFile  = errors.New("failed to read config")
	ErrFailedToParseFile = errors.New("failed to parse config")
)

type Rule struct {
	Key     string `yaml:"key"`
	Exact   bool   `yaml:"exact"`
	Prefix  bool   `yaml:"prefix"`
	Suffix  bool   `yaml:"suffix"`
	Include bool   `yaml:"include"`
}

type Rules struct {
	Rules []Rule `yaml:"rules"`
}

type Config struct {
	Rules           []Rule   `yaml:"rules"`
	IgnoreEnvFiles  []string `yaml:"ignoreEnvFiles"`
	IgnoreYAMLFiles []string `yaml:"ignoreYamlFiles"`
}

func Load(configFilePath string) (Config, error) {
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, ErrFailedToReadFile
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, ErrFailedToParseFile
	}

	return config, nil
}

func GetDefaultConfig() Config {
	var config Config

	config.IgnoreEnvFiles = make([]string, 0)
	config.IgnoreEnvFiles = append(config.IgnoreEnvFiles, ".env.example")
	config.IgnoreEnvFiles = append(config.IgnoreEnvFiles, ".env.sample")

	config.IgnoreYAMLFiles = make([]string, 0)
	config.IgnoreYAMLFiles = append(config.IgnoreYAMLFiles, ".example.yml")

	config.Rules = append(config.Rules, addDefaultRule("PRIVATE_KEY"))
	config.Rules = append(config.Rules, addDefaultRule("MNEMONIC"))
	config.Rules = append(config.Rules, addDefaultRule("API_KEY"))
	config.Rules = append(config.Rules, addDefaultRule("TOKEN"))
	config.Rules = append(config.Rules, addDefaultRule("PASSWORD"))
	config.Rules = append(config.Rules, addDefaultRule("PASS"))
	config.Rules = append(config.Rules, addDefaultRule("SECRET"))

	return config
}

func addDefaultRule(key string) Rule {
	return Rule{
		Key:     key,
		Exact:   true,
		Prefix:  true,
		Suffix:  true,
		Include: true,
	}
}
