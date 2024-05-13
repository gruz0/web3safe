package yamlanalyzer

import "flag"

type Flags struct {
	PathToScan     string
	ConfigFilePath string
}

func ParseFlags() Flags {
	var flags Flags

	flag.StringVar(&flags.PathToScan, "path", ".", "Path to scan")
	flag.StringVar(&flags.ConfigFilePath, "config", "", "Path to the configuration file")
	flag.Parse()

	return flags
}
