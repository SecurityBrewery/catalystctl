package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

type URLConfig struct {
	// CatalystURL is the URL of the Catalyst server.
	URL string `yaml:"url"`
}

type TokenConfig struct {
	// Token is the authentication token to use.
	Token string `yaml:"token"`
}

func URL(u *url.URL) (*url.URL, error) {
	if u != nil {
		return u, nil
	}

	catalystURL := os.Getenv("CATALYST_URL")
	if catalystURL != "" {
		u, err := url.Parse(catalystURL)
		if err == nil {
			return u, nil
		} else {
			log.Println("failed to parse CATALYST_URL:", err)
		}
	}

	path := configFilePath()
	if _, err := os.Stat(path); err != nil {
		return nil, errors.New("please set `--url`, the environment variable `CATALYST_URL`, or set the `url: <url>` in the config file at " + path)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cli := &URLConfig{}
	if err := yaml.Unmarshal(b, cli); err != nil {
		return nil, err
	}

	if cli.URL != "" {
		u, err := url.Parse(cli.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse url from config file: %w", err)
		}
		return u, nil
	}

	return nil, errors.New("please set `--url`, the environment variable `CATALYST_URL`, or set the `url: <url>` in the config file at " + path)
}

func Token(t string) (string, error) {
	if t != "" {
		return t, nil
	}

	token := os.Getenv("CATALYST_TOKEN")
	if token != "" {
		return token, nil
	}

	path := configFilePath()
	if _, err := os.Stat(path); err != nil {
		return "", errors.New("please set `--token`, the environment variable `CATALYST_TOKEN`, or set the `token: <token>` in the config file at " + path)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	cli := &TokenConfig{}
	if err := yaml.Unmarshal(b, cli); err != nil {
		return "", err
	}

	if cli.Token == "" {
		return "", errors.New("please set `--token`, the environment variable `CATALYST_TOKEN`, or set the `token: <token>` in the config file at " + path)
	}

	return cli.Token, nil
}

func CreateFile(url, token string) error {
	b, err := yaml.Marshal(map[string]string{
		"url":   url,
		"token": token,
	})
	if err != nil {
		return err
	}

	path := configFilePath()

	_ = os.Mkdir(filepath.Dir(path), 0755)
	return os.WriteFile(path, b, 0600)
}

func configFilePath() string {
	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".catalyst", "config.yaml")
}
