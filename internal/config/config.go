package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

const (
	appName             = "web3safe"
	configDirPermission = 0o755
	configPermission    = 0o600
)

var (
	ErrFailedToReadFile   = errors.New("failed to read file")
	ErrFailedToParseFile  = errors.New("failed to parse file")
	ErrFailedToEncodeFile = errors.New("failed to encode file")
	ErrFailedToWriteFile  = errors.New("failed to write file")
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
	IgnoreFiles     []string `yaml:"ignoreFiles"`
	IgnoreYAMLFiles []string `yaml:"ignoreYamlFiles"`
}

func LoadConfig(configFilePath string) (Config, error) {
	var config Config

	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, errors.Join(ErrFailedToReadFile, err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, errors.Join(ErrFailedToParseFile, err)
	}

	return config, nil
}

func GetDefaultConfig() Config {
	var config Config

	config.IgnoreFiles = make([]string, 0)
	config.IgnoreFiles = append(config.IgnoreFiles, ".env.example")
	config.IgnoreFiles = append(config.IgnoreFiles, ".env.sample")

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

func GenerateConfig(filePath string) error {
	yamlData, err := yaml.Marshal(GetDefaultConfig())
	if err != nil {
		return errors.Join(ErrFailedToEncodeFile, err)
	}

	err = os.WriteFile(filePath, yamlData, configPermission)
	if err != nil {
		return errors.Join(ErrFailedToWriteFile, err)
	}

	return nil
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

func GetDefaultConfigFilePath() string {
	var configDir string

	switch osType := runtime.GOOS; osType {
	case "windows":
		appData := os.Getenv("APPDATA")
		configDir = filepath.Join(appData, appName)
	case "darwin", "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("Unable to get home directory")
		}

		configDir = filepath.Join(homeDir, ".config", appName)
	default:
		panic("Unsupported operating system")
	}

	if err := os.MkdirAll(configDir, configDirPermission); err != nil {
		panic(err)
	}

	return filepath.Join(configDir, "config.yaml")
}
