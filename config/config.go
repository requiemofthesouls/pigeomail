package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port int64 `yaml:"port"`
}

// Load parses the YAML input into a Config.
func Load(input []byte) (cfg *Config, err error) {
	if err = yaml.UnmarshalStrict(input, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadFile parses the given YAML file into a Config
func LoadFile(filename string) (cfg *Config, err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err != nil {
		return nil, err
	}

	return Load(content)

}
