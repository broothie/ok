package cli

import (
	"os"

	"github.com/bmatcuk/doublestar"
	"github.com/broothie/ok/logger"
	"github.com/kelseyhightower/envconfig"
	"github.com/pelletier/go-toml"
)

const configFileGlob = ".ok*.toml"

type Options struct {
	Debug        bool     `toml:"debug" envconfig:"debug"`
	Help         bool     `toml:"-"`
	Version      bool     `toml:"-"`
	Init         string   `toml:"-"`
	ListTools    bool     `toml:"-"`
	Watches      []string `toml:"-"`
	Halt         bool     `toml:"-"`
	SkipTools    []string `toml:"skip" envconfig:"skip"`
	ToolPriority []string `toml:"tool_priority" envconfig:"tool_priority"`
	TaskName     string   `toml:"-"`
}

func ReadInNewConfig() Options {
	filenames, err := doublestar.Glob(configFileGlob)
	if err != nil {
		logger.Ok.Printf("failed to glob for '%s'", configFileGlob)
	}

	var config Options
	if err := envconfig.Process("ok", &config); err != nil {
		logger.Ok.Printf("failed to read config from env: %v", err)
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
