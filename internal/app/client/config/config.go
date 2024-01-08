package config

import (
	"encoding/json"
	"fmt"

	"github.com/alecthomas/kingpin/v2"
)

const (
	defaultDebugLevel = "0"
	defaultAddress    = "127.0.0.1:3200"
)

type Config struct {
	Debug   int
	Version bool
	Address string
}

func New() *Config {
	var c Config
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Flag("debug", "debug modes: 0..2").
		IntVar(&c.Debug)
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
