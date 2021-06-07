package config

import (
	"os"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/broothie/ok/logger"
	"github.com/kelseyhightower/envconfig"
	"github.com/pelletier/go-toml"
)

const configFileGlob = ".ok*.toml"

type Config struct {
	Debug        bool          `toml:"debug" envconfig:"debug"`
	Timeout      time.Duration `toml:"timeout" envconfig:"timeout" default:"1ms"`
	SkipTools    []string      `toml:"skip" envconfig:"skip"`
	ToolPriority []string      `toml:"tool_priority" envconfig:"tool_priority"`
}

func ReadConfigAndEnv(v interface{}) {
	ReadInEnv(v)
	ReadInConfig(v)
}

func ReadInEnv(v interface{}) {
	if err := envconfig.Process("ok", v); err != nil {
		logger.Ok.Printf("failed to read config from env: %v", err)
	}
}

func ReadInConfig(v interface{}) {
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

		if err := toml.NewDecoder(file).Decode(v); err != nil {
			logger.Ok.Printf("failed to decode config file '%s': %v", filename, err)
			continue
		}
	}
}
