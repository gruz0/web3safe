package dotenvanalyzer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gruz0/web3safe/internal/config"
	"github.com/gruz0/web3safe/internal/utils"
)

const (
	envVariablePartsCount = 2
	reportMessageFormat   = "%s:%d: found sensitive variable %s"
)

var ErrEnvVariableHasInvalidFormat = errors.New("variable has invalid ENV format")

type DotEnvAnalyzer struct {
	config        config.Config
	report        []string
	filesToIgnore map[string]bool
}

func NewDotEnvAnalyzer(config config.Config) *DotEnvAnalyzer {
	filesToIgnore := make(map[string]bool)
	for _, fileToIgnore := range config.IgnoreEnvFiles {
		filesToIgnore[fileToIgnore] = true
	}

	return &DotEnvAnalyzer{
		config:        config,
		report:        make([]string, 0),
		filesToIgnore: filesToIgnore,
	}
}

func (s *DotEnvAnalyzer) ScanDirectory(directoryPath string, recursive bool) error {
	if recursive {
		return s.scanDirectoryRecursively(directoryPath)
	}

	return s.scanDirectory(directoryPath)
}

func (s *DotEnvAnalyzer) ScanFile(filePath string) error {
	filename := filepath.Base(filePath)

	if s.isFileSkippable(filename) {
		return nil
	}

	return s.analyzeVariables(filePath)
}

func (s *DotEnvAnalyzer) Report() []string {
	return s.report
}

func (s *DotEnvAnalyzer) scanDirectoryRecursively(directoryPath string) error {
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

		if err := s.analyzeVariables(path); err != nil {
			return fmt.Errorf("failed to analyze file: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory recursively: %w", err)
	}

	return nil
}

func (s *DotEnvAnalyzer) scanDirectory(directoryPath string) error {
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

		if err := s.analyzeVariables(filepath.Join(directoryPath, entry.Name())); err != nil {
			return fmt.Errorf("failed to analyze file: %w", err)
		}
	}

	return nil
}

func (s *DotEnvAnalyzer) isFileSkippable(fileName string) bool {
	if !strings.HasPrefix(fileName, ".env") {
		return true
	}

	if s.filesToIgnore[fileName] {
		return true
	}

	return false
}

//nolint:funlen,gocognit,cyclop
func (s *DotEnvAnalyzer) analyzeVariables(filePath string) error {
	lines, err := utils.GetFileContent(filePath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for idx, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimPrefix(line, "export ")

		parts := strings.SplitN(line, "=", envVariablePartsCount)
		if len(parts) != envVariablePartsCount {
			return fmt.Errorf("%w: %s", ErrEnvVariableHasInvalidFormat, line)
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
