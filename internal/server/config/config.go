package config

import (
	"github.com/alecthomas/kingpin/v2"
)

const (
	defaultLogLevel = "0"
)

type Config struct {
	LogLevel int8
	Debug    bool
}

func New() *Config {
	var c Config
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Flag("LogLevel", "-1..2, where -1=Debug 0=Info 1=Warning 2=Error").
		Short('l').
		Envar("LOGLEVEL").
		Default(defaultLogLevel).
		Int8Var(&c.LogLevel)
	kingpin.Flag("debug", "enable debug mode").
		BoolVar(&c.Debug)
	kingpin.Parse()
	return &c
}
