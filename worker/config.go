package worker

import (
	"github.com/BurntSushi/toml"
)

// config represents worker config options
type config struct {
	Global globalConfig `toml:"global"`
	Server serverConfig `toml:"server"`
	Client clientConfig `toml:"client"`
}

type globalConfig struct {
	Name   string `toml:"name"`
	LogDir string `toml:"log_dir"`
	JobDir string `toml:"job_dir"`
}

type serverConfig struct {
	Hostname string `toml:"hostname"`
	Addr     string `toml:"listen_addr"`
	Port     int    `toml:"listen_port"`

	CACert  string `toml:"ca_cert"`
	SSLCert string `toml:"ssl_cert"`
	SSLKey  string `toml:"ssl_key"`
}

type clientConfig struct {
	CACert  string `toml:"ca_cert"`
	SSLCert string `toml:"ssl_cert"`
	SSLKey  string `toml:"ssl_key"`
}

// LoadConfig loads configuration
func LoadConfig(cfgFile string) (*config, error) {
	cfg := new(config)
	if _, err := toml.DecodeFile(cfgFile, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
