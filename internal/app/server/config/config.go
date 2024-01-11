package config

import (
	"encoding/json"
	"fmt"

	"github.com/alecthomas/kingpin/v2"
)

const (
	defaultLogLevel = "0"
	defaultDbUri    = "mongodb://mongo:27017"
	defaultAddress  = "0.0.0.0:3200"
)

type Config struct {
	LogLevel    int
	Debug       bool
	Version     bool
	DbUri       string
	Address     string
	TokenKey    string
	TLSCertFile string
	TLSKeyFile  string
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
	kingpin.Flag("version", "print version and exit").Short('v').BoolVar(&c.Version)
	kingpin.Flag("db-uri", "database connection string").
		Short('d').
		Envar("DB_URI").
		Default(defaultDbUri).
		StringVar(&c.DbUri)
	kingpin.Flag("address", "server listen address host:port").
		Short('a').
		Envar("LISTEN_ADDRESS").
		Default(defaultAddress).
		StringVar(&c.Address)
	kingpin.Flag("token-key", "key to sign tokens").
		Short('k').
		Envar("TOKEN_KEY").
		StringVar(&c.TokenKey)
	kingpin.Flag(
		"tls-cert-file",
		"path to server tls certificate file. This flag mandatory enables tls. Certificate file should contain full certificate chain, including intermediate CA certificates (if any).",
	).Envar("TLS_CERT_FILE").StringVar(&c.TLSCertFile)
	kingpin.Flag(
		"tls-key-file",
		"path to server tls certificate key file",
	).Envar("TLS_KEY_FILE").StringVar(&c.TLSKeyFile)
	kingpin.Parse()
	return &c
}

func (c Config) Print() {
	b, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println("Configuration:")
	fmt.Println(string(b))
}
