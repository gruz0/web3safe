package envanalyzer

import (
	"errors"
	"os"
	"strings"

	"github.com/gruz0/web3safe/internal/config"
)

const envVariablePartsCount = 2

var (
	ErrFailedToAnalyzeEnvVariables = errors.New("failed to analyze env variables")
	ErrEnvVariableHasInvalidFormat = errors.New("variable has invalid ENV format")
)

type EnvAnalyzer struct {
	config config.Config
	report []string
}

func NewEnvAnalyzer(config config.Config) *EnvAnalyzer {
	return &EnvAnalyzer{
		config: config,
		report: make([]string, 0),
	}
}

func (s *EnvAnalyzer) Run() error {
	variables := os.Environ()

	if err := s.analyzeVariables(variables); err != nil {
		return errors.Join(ErrFailedToAnalyzeEnvVariables, err)
	}

	return nil
}

func (s *EnvAnalyzer) Report() []string {
	return s.report
}

//nolint:funlen,gocognit,cyclop
func (s *EnvAnalyzer) analyzeVariables(variables []string) error {
	for _, line := range variables {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

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

		matches := make([]string, 0)
		encountered := make(map[string]bool)

		for _, rule := range s.config.Rules {
			lowercasedKey := strings.ToLower(rule.Key)

			if rule.Exact && lowercasedVariable == lowercasedKey {
				if encountered[variable] {
					break
				}

				encountered[variable] = true

				matches = append(matches, variable)
			}

			if rule.Prefix && strings.HasPrefix(lowercasedVariable, lowercasedKey+"_") {
				if encountered[variable] {
					break
				}

				encountered[variable] = true

				matches = append(matches, variable)
			}

			if rule.Suffix && strings.HasSuffix(lowercasedVariable, "_"+lowercasedKey) {
				if encountered[variable] {
					break
				}

				encountered[variable] = true

				matches = append(matches, variable)
			}

			if rule.Include && strings.Contains(lowercasedVariable, "_"+lowercasedKey+"_") {
				if encountered[variable] {
					break
				}

				encountered[variable] = true

				matches = append(matches, variable)
			}
		}

		for _, match := range matches {
			s.report = append(s.report, "Shell ENV has a sensitive variable: "+match)
		}
	}

	return nil
}
