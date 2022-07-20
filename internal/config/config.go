package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "merkle.yml"

type Config struct {
	ClientTimeout   time.Duration `yaml:"client_timeout"`
	FirstPieceIndex int           `yaml:"first_piece_index"`
	LastPieceIndex  int           `yaml:"last_piece_index"`
	IconsHash       string        `yaml:"icons_hash"`
	ServerURL       string        `yaml:"server_url"`
}

func Get(filePath string) (*Config, error) {
	c := Config{}
	c.setDefaults()
	err := load(filePath, &c)
	if err != nil {
		return nil, err
	}
	err = validate(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func load(filePath string, cfg *Config) error {
	if filePath == "" {
		filePath = DefaultConfigPath
	}
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) setDefaults() {
	c.ClientTimeout = 5 * time.Second
}

func validate(cfg *Config) error {
	if cfg.ClientTimeout == 0 {
		return errors.New("config: client timeout not defined")
	}
	if cfg.LastPieceIndex == 0 {
		return errors.New("config: last piece index not defined")
	}
	if cfg.IconsHash == "" {
		return errors.New("config: icons hash not defined")
	}
	if cfg.ServerURL == "" {
		return errors.New("config: server URL not defined")
	}
	return nil
}
