package envanalyzer

import "flag"

type Flags struct {
	ConfigFilePath string
}

func ParseFlags() Flags {
	var flags Flags

	flag.StringVar(&flags.ConfigFilePath, "config", "", "Path to the configuration file")
	flag.Parse()

	return flags
}
