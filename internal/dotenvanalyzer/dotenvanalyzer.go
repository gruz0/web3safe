package dotenvanalyzer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gruz0/web3safe/internal/config"
)

const (
	envVariablePartsCount = 2
	reportMessageFormat   = "%s:%d: found sensitive variable %s"
)

var (
	ErrFailedToAccessPath          = errors.New("failed to access path")
	ErrFailedToAnalyzeVariables    = errors.New("failed to analyze variables")
	ErrEnvVariableHasInvalidFormat = errors.New("variable has invalid ENV format")
)

type DotEnvAnalyzer struct {
	pathToScan string
	config     config.Config
	report     []string
}

func NewDotEnvAnalyzer(pathToScan string, config config.Config) *DotEnvAnalyzer {
	return &DotEnvAnalyzer{
		pathToScan: pathToScan,
		config:     config,
		report:     make([]string, 0),
	}
}

func (s *DotEnvAnalyzer) Run() error {
	filesToIgnore := make(map[string]bool)
	for _, fileToIgnore := range s.config.IgnoreEnvFiles {
		filesToIgnore[fileToIgnore] = true
	}

	err := filepath.Walk(s.pathToScan, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Join(ErrFailedToAccessPath, err)
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		if !s.hasEnvPrefix(info.Name()) {
			return nil
		}

		if filesToIgnore[info.Name()] {
			return nil
		}

		if err := s.analyzeVariables(path); err != nil {
			return errors.Join(ErrFailedToAnalyzeVariables, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	return nil
}

func (s *DotEnvAnalyzer) Report() []string {
	return s.report
}

func (s *DotEnvAnalyzer) hasEnvPrefix(name string) bool {
	return strings.HasPrefix(name, ".env")
}

//nolint:funlen,gocognit,cyclop
func (s *DotEnvAnalyzer) analyzeVariables(filePath string) error {
	lines, err := getFileContent(filePath)
	if err != nil {
		return err
	}

	for idx, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimPrefix(line, "export ")

		parts := strings.SplitN(line, "=", envVariablePartsCount)
		if len(parts) != envVariablePartsCount {
			return ErrEnvVariableHasInvalidFormat
		}

		// If value is empty, just skip
		if parts[1] == "" {
			continue
		}

		variable := strings.TrimSpace(parts[0])
		lowercasedVariable := strings.ToLower(variable)

		type match struct {
			Line     int
			Variable string
		}

		matches := make([]match, 0)
		encountered := make(map[string]bool)

		for _, rule := range s.config.Rules {
			lowercasedKey := strings.ToLower(rule.Key)

			if rule.Exact && lowercasedVariable == lowercasedKey {
				if encountered[variable] {
					break
				}

				encountered[variable] = true

				matches = append(matches, match{Line: idx + 1, Variable: variable})
			}

			if rule.Prefix && strings.HasPrefix(lowercasedVariable, lowercasedKey+"_") {
				if encountered[variable] {
					break
				}

				encountered[variable] = true

				matches = append(matches, match{Line: idx + 1, Variable: variable})
			}

			if rule.Suffix && strings.HasSuffix(lowercasedVariable, "_"+lowercasedKey) {
				if encountered[variable] {
					break
				}

				encountered[variable] = true

				matches = append(matches, match{Line: idx + 1, Variable: variable})
			}

			if rule.Include && strings.Contains(lowercasedVariable, "_"+lowercasedKey+"_") {
				if encountered[variable] {
					break
				}

				encountered[variable] = true

				matches = append(matches, match{Line: idx + 1, Variable: variable})
			}
		}

		for _, match := range matches {
			s.report = append(s.report, fmt.Sprintf(reportMessageFormat, filePath, match.Line, match.Variable))
		}
	}

	return nil
}

func getFileContent(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return strings.Split(string(content), "\n"), nil
}
