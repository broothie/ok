package cli

import (
	"os"

	"github.com/bmatcuk/doublestar"
	"github.com/broothie/ok/logger"
	"github.com/kelseyhightower/envconfig"
	"github.com/pelletier/go-toml"
)

const configFileGlob = ".ok*.toml"

type Config struct {
	Debug        bool     `toml:"debug" envconfig:"debug"`
	SkipTools    []string `toml:"skip" envconfig:"skip"`
	ToolPriority []string `toml:"tool_priority" envconfig:"tool_priority"`
}

type Options struct {
	Config
	Help      bool     `toml:"-"`
	Version   bool     `toml:"-"`
	Init      string   `toml:"-"`
	ListTools bool     `toml:"-"`
	Watches   []string `toml:"-"`
}

func ReadInConfig() Config {
	var config Config
	if err := envconfig.Process("ok", &config); err != nil {
		logger.Ok.Printf("failed to read config from env: %v", err)
	}

	filenames, err := doublestar.Glob(configFileGlob)
	if err != nil {
		logger.Ok.Printf("failed to glob for '%s'", configFileGlob)
	}

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			logger.Ok.Printf("failed to read config file '%s': %v", filename, err)
			continue
		}

		if err := toml.NewDecoder(file).Decode(&config); err != nil {
			logger.Ok.Printf("failed to decode config file '%s': %v", filename, err)
			continue
		}
	}

	return config
}
