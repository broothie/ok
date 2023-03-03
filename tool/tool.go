package tool

import (
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
)

type NewFunc func() Tool

type Tool interface {
	Name() string
	Config() *Config
	Init() error
	ProcessFile(path string) ([]task.Task, error)
}

type Config map[string]string

func (c Config) Entries() map[string]string {
	return c
}

func (c Config) Set(key, value string) {
	c[key] = value
}

func (c Config) Get(key string) string {
	return c[key]
}

func (c Config) Executable() string {
	return c.Get("executable")
}

func (c Config) Filenames() []string {
	return util.SplitCommaList(c.Get("filenames"))
}

func (c Config) Extensions() []string {
	return util.SplitCommaList(c.Get("extensions"))
}
