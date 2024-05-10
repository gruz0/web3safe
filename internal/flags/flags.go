package flags

import "flag"

type Flags struct {
	GenerateConfig bool
}

func ParseFlags() Flags {
	var flags Flags

	flag.BoolVar(&flags.GenerateConfig, "generateConfig", false, "Generate a new configuration file")
	flag.Parse()

	return flags
}
