package config

import (
	"encoding/json"
	"fmt"

	"github.com/alecthomas/kingpin/v2"
)

const (
	defaultLogLevel = "0"
	defaultAddress  = "127.0.0.1:3200"
)

type Config struct {
	LogLevel int
	Debug    bool
	Version  bool
	Address  string
}

func New() *Config {
	var c Config
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Flag("LogLevel", "-1..2, where -1=Debug 0=Info 1=Warning 2=Error").
		Short('l').
		Envar("LOGLEVEL").
		Default(defaultLogLevel).
		IntVar(&c.LogLevel)
	kingpin.Flag("debug", "enable debug mode").
		BoolVar(&c.Debug)
	kingpin.Flag("version", "print version").Short('v').BoolVar(&c.Version)
	kingpin.Flag("address", "server address host:port").
		Short('a').
		Envar("LISTEN_ADDRESS").
		Default(defaultAddress).
		StringVar(&c.Address)
	kingpin.Parse()
	return &c
}

func (c Config) Print() {
	b, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println("Configuration:")
	fmt.Println(string(b))
}