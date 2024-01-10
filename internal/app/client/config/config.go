package config

import (
	"encoding/json"
	"fmt"

	"github.com/alecthomas/kingpin/v2"
)

const (
	defaultAddress = "127.0.0.1:3200"
)

type Config struct {
	Logfile  string
	Log      bool
	Version  bool
	Address  string
	UseMouse bool
}

func New() *Config {
	var c Config
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Flag("log", "enable logging").
		Short('l').
		BoolVar(&c.Log)
	kingpin.Flag("logfile", "log file name").
		Envar("LOGFILE").
		StringVar(&c.Logfile)
	kingpin.Flag("version", "print version").Short('v').BoolVar(&c.Version)
	kingpin.Flag("address", "server address host:port").
		Short('a').
		Envar("LISTEN_ADDRESS").
		Default(defaultAddress).
		StringVar(&c.Address)
	kingpin.Flag("mouse", "enable mouse support (may be unstable)").
		Short('m').
		BoolVar(&c.UseMouse)
	kingpin.Parse()
	return &c
}

func (c Config) Print() {
	b, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println("Configuration:")
	fmt.Println(string(b))
}
