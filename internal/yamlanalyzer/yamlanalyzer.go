package yamlanalyzer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gruz0/web3safe/internal/config"
	"gopkg.in/yaml.v2"
)

const (
	reportMessageFormat = "%s: found sensitive key \"%s\" in %s"
)

var (
	ErrFailedToAccessPath       = errors.New("failed to access path")
	ErrFailedToParseFile        = errors.New("failed to parse file")
	ErrFailedToAnalyzeVariables = errors.New("failed to analyze variables")
)

type YamlAnalyzer struct {
	pathToScan string
	config     config.Config
	report     []string
}

type match struct {
	Path     string
	Variable string
}

func NewYamlAnalyzer(pathToScan string, config config.Config) *YamlAnalyzer {
	return &YamlAnalyzer{
		pathToScan: pathToScan,
		config:     config,
		report:     make([]string, 0),
	}
}

func (s *YamlAnalyzer) Run() error {
	err := filepath.Walk(s.pathToScan, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Join(ErrFailedToAccessPath, err)
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		if !s.isYAML(info.Name()) {
			return nil
		}

		if info.Size() == 0 {
			return nil
		}

		for _, fileToIgnore := range s.config.IgnoreYAMLFiles {
			if strings.Contains(info.Name(), fileToIgnore) {
				return nil
			}
		}

		if err := s.analyzeKeys(path); err != nil {
			return errors.Join(ErrFailedToAnalyzeVariables, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	return nil
}

func (s *YamlAnalyzer) Report() []string {
	return s.report
}

func (s *YamlAnalyzer) isYAML(name string) bool {
	return strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")
}

func (s *YamlAnalyzer) analyzeKeys(filePath string) error {
	yamlData, err := getFileContent(filePath)
	if err != nil {
		return err
	}

	var data interface{}
	if err := yaml.Unmarshal(yamlData, &data); err != nil {
		return errors.Join(ErrFailedToParseFile, err)
	}

	keys := traverse(data)

	if len(keys) == 0 {
		return nil
	}

	for _, path := range keys {
		parts := strings.Split(path, ".")

		for _, part := range parts {
			matches := s.findMatches(part, path)

			for _, match := range matches {
				s.report = append(s.report, fmt.Sprintf(reportMessageFormat, filePath, match.Variable, match.Path))
			}
		}
	}

	return nil
}

func getFileContent(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return content, nil
}

// Thanks https://github.com/egibs/deepwalk/blob/main/pkg/traverse/traverse.go!
func traverse(value interface{}) []string {
	var processQueue func(value reflect.Value, path string)

	keys := make([]string, 0)

	processQueue = func(value reflect.Value, path string) {
		// Handle pointers and interface{} by obtaining the value they point to
		for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
			if value.IsNil() {
				return
			}

			value = value.Elem()
		}

		if value.Kind() == reflect.Map {
			for _, key := range value.MapKeys() {
				val := value.MapIndex(key)

				fullPath := fmt.Sprintf("%s.%v", path, key.Interface())

				keys = append(keys, fullPath)

				processQueue(val, fullPath)
			}

			return
		}

		if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
			for i := range make([]struct{}, value.Len()) {
				processQueue(value.Index(i), fmt.Sprintf("%s[%d]", path, i))
			}

			return
		}
	}

	// Convert the interface{} to reflect.Value before processing
	initialValue := reflect.ValueOf(value)
	processQueue(initialValue, "")

	return keys
}

func (s *YamlAnalyzer) findMatches(part, path string) []match {
	lowercasedPart := strings.ToLower(part)

	matches := make([]match, 0)
	encountered := make(map[string]bool)

	ruleFunc := func(condition bool, key string) {
		if !condition {
			return
		}

		if encountered[key] {
			return
		}

		encountered[key] = true

		matches = append(matches, match{Path: path, Variable: part})
	}

	for _, rule := range s.config.Rules {
		lowercasedRuleKey := strings.ToLower(rule.Key)

		ruleFunc(rule.Exact && lowercasedPart == lowercasedRuleKey, lowercasedPart)
		ruleFunc(rule.Prefix && strings.HasPrefix(lowercasedPart, lowercasedRuleKey+"_"), lowercasedPart)
		ruleFunc(rule.Suffix && strings.HasSuffix(lowercasedPart, "_"+lowercasedRuleKey), lowercasedPart)
		ruleFunc(rule.Include && strings.Contains(lowercasedPart, "_"+lowercasedRuleKey+"_"), lowercasedPart)
	}

	return matches
}
