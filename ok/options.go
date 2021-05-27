package ok

import (
	"os"

	"github.com/bmatcuk/doublestar"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Options struct {
	Debug               bool
	Help                bool
	Version             bool
	Init                string
	ListTools           bool
	Watches             []string
	Stop                bool
	TaskName            string
	ConfigurableOptions ConfigMap
}

func NewOptionsFromEnvironment() (Options, error) {
	filenames, err := doublestar.Glob("*.ok*.toml")
	if err != nil {
		return Options{}, errors.Wrap(err, "failed to glob for '*.ok*.toml'")
	}

	var config ConfigMap
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			Logger.Printf("failed to read config file '%s': %v", filename, err)
			continue
		}

		if err := toml.NewDecoder(file).Decode(&config); err != nil {
			Logger.Printf("failed to decode config file '%s': %v", filename, err)
			continue
		}
	}

	return Options{Debug: config.Debug(), ConfigurableOptions: config}, nil
}

type Config struct {
	Debug bool
}

func NewConfig() (Config, error) {
	cfg := viper.New()
	cfg.SetConfigName(".ok")
	if err := cfg.ReadInConfig(); err != nil {
		return Config{}, errors.Wrap(err, "")
	}

	return Config{}, nil
}
