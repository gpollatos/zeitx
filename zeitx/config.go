package zeitx

import (
	"errors"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// HTTPServerConfig is a straightforward struct for the http config values
type HTTPServerConfig struct {
	ListenAddr string `yaml:"listen"`
}

// Config is a straightforward struct for the config.yml file
type Config struct {
	HTTPServer HTTPServerConfig `yaml:"httpsrv"`
}

// NewConfig creates a new Config struct from the provided fileName
func NewConfig(fileName string) (*Config, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, err
	}
	if cfg.HTTPServer.ListenAddr == "" {
		return nil, errors.New("listen address cannot be empty")
	}
	return cfg, nil
}
