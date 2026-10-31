package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

type ConfigServer struct {
	interval time.Duration

	Port          int    `yaml:"port"`
	FetchInterval string `yaml:"fetch_interval"`
	GitHubToken   string `yaml:"github_token"`
}

type Config struct {
	Server  ConfigServer `yaml:"server"`
	Targets []Target     `yaml:"targets"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

func (c *Config) Init() error {
	// check if port is valid
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server.port %d", c.Server.Port)
	}

	// default timeout
	if c.Server.FetchInterval == "" {
		return errors.New("missing server.fetch_interval")
	}

	// parse interval
	interval, err := ParseDuration(c.Server.FetchInterval)
	if err != nil {
		return err
	}

	c.Server.interval = interval

	if len(c.Targets) == 0 {
		log.Warnln("no targets configured")
	}

	for i, target := range c.Targets {
		err := target.Validate()
		if err != nil {
			return fmt.Errorf("targets[%d]: %s", i, err)
		}
	}

	return nil
}

func (c *Config) Store() error {
	var (
		buffer   bytes.Buffer
		comments = yaml.CommentMap{
			"$.targets": {yaml.HeadComment("")},

			"$.server.port":           {yaml.HeadComment(" port to serve gitmon on (required; default 9229)")},
			"$.server.fetch_interval": {yaml.HeadComment(" fetch interval (required; default 30m)")},
			"$.server.github_token":   {yaml.HeadComment(" github api token (optional)")},
		}
	)

	err := yaml.NewEncoder(&buffer, yaml.WithComment(comments)).Encode(c)
	if err != nil {
		return err
	}

	body := bytes.ReplaceAll(buffer.Bytes(), []byte("#\n"), []byte("\n"))

	err = os.WriteFile("config.yml", body, 0600)
	if err != nil {
		return err
	}

	return os.Chmod("config.yml", 0600)
}

func LoadConfig() (*Config, error) {
	// defaults
	cfg := &Config{
		Server: ConfigServer{
			Port:          9229,
			FetchInterval: "30m",
		},
	}

	file, err := os.Open("config.yml")
	if os.IsNotExist(err) {
		err = cfg.Store()
		if err != nil {
			return nil, err
		}
	} else {
		if err != nil {
			return nil, err
		}

		defer file.Close()

		err = yaml.NewDecoder(file).Decode(cfg)
		if err != nil {
			return nil, err
		}
	}

	err = cfg.Init()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
