package yamlanalyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gruz0/web3safe/internal/config"
	"gopkg.in/yaml.v2"
)

const (
	reportMessageFormat = "%s: found sensitive key %q in %s"
)

type YamlAnalyzer struct {
	config config.Config
	report []string
}

type match struct {
	Path     string
	Variable string
}

func NewYamlAnalyzer(config config.Config) *YamlAnalyzer {
	return &YamlAnalyzer{
		config: config,
		report: make([]string, 0),
	}
}

func (s *YamlAnalyzer) ScanDirectory(directoryPath string, recursive bool) error {
	if recursive {
		return s.scanDirectoryRecursively(directoryPath)
	}

	return s.scanDirectory(directoryPath)
}

func (s *YamlAnalyzer) ScanFile(filePath string) error {
	filename := filepath.Base(filePath)

	if s.isFileSkippable(filename) {
		return nil
	}

	return s.analyzeKeys(filePath)
}

func (s *YamlAnalyzer) Report() []string {
	return s.report
}

func (s *YamlAnalyzer) scanDirectoryRecursively(directoryPath string) error {
	err := filepath.WalkDir(directoryPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access directory: %w", err)
		}

		if entry.IsDir() {
			return nil
		}

		if s.isFileSkippable(entry.Name()) {
			return nil
		}

		if err := s.analyzeKeys(path); err != nil {
			return fmt.Errorf("failed to analyze file: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory recursively: %w", err)
	}

	return nil
}

func (s *YamlAnalyzer) scanDirectory(directoryPath string) error {
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		if s.isFileSkippable(entry.Name()) {
			continue
		}

		if err := s.analyzeKeys(filepath.Join(directoryPath, entry.Name())); err != nil {
			return fmt.Errorf("failed to analyze file: %w", err)
		}
	}

	return nil
}

func (s *YamlAnalyzer) isFileSkippable(fileName string) bool {
	if !strings.HasSuffix(fileName, ".yaml") && !strings.HasSuffix(fileName, ".yml") {
		return true
	}

	for _, fileToIgnore := range s.config.IgnoreYAMLFiles {
		if strings.Contains(fileName, fileToIgnore) {
			return true
		}
	}

	return false
}

func (s *YamlAnalyzer) analyzeKeys(filePath string) error {
	yamlData, err := getFileContent(filePath)
	if err != nil {
		return err
	}

	var data interface{}
	if err := yaml.Unmarshal(yamlData, &data); err != nil {
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
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
