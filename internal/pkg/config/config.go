package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	DEFENSE_HOST                  string
	DEFENSE_PORT                  string
	DEFENSE_SCHEME                string
	DEFENSE_INSECURE_SKIP_VERIFY  bool
	DEFENSE_USER                  string
	DEFENSE_PASS                  string
	HTTP_ADDR                     string
	HTTP_API_KEY                  string
}

func (config *Config) Load() error {
	config.DEFENSE_HOST = os.Getenv("DEFENSE_HOST")
	config.DEFENSE_PORT = os.Getenv("DEFENSE_PORT")
	config.DEFENSE_SCHEME = os.Getenv("DEFENSE_SCHEME")
	config.DEFENSE_INSECURE_SKIP_VERIFY, _ = strconv.ParseBool(os.Getenv("DEFENSE_INSECURE_SKIP_VERIFY"))
	config.DEFENSE_USER = os.Getenv("DEFENSE_USERNAME")
	config.DEFENSE_PASS = os.Getenv("DEFENSE_PASSWORD")
	config.HTTP_ADDR = os.Getenv("HTTP_ADDR")
	config.HTTP_API_KEY = os.Getenv("HTTP_API_KEY")

	if config.DEFENSE_SCHEME == "" {
		config.DEFENSE_SCHEME = "https"
	}

	if config.HTTP_ADDR == "" {
		config.HTTP_ADDR = ":8080"
	}

	if config.DEFENSE_HOST == "" {
		return errors.New("DEFENSE_HOST cant'n be null")
	}

	if config.DEFENSE_PORT == "" {
		return errors.New("DEFENSE_PORT cant'n be null")
	}

	if config.DEFENSE_USER == "" {
		return errors.New("DEFENSE_USERNAME cant'n be null")
	}

	if config.DEFENSE_PASS == "" {
		return errors.New("DEFENSE_PASSWORD cant'n be null")
	}

	if config.HTTP_API_KEY == "" {
		return errors.New("HTTP_API_KEY cant'n be null")
	}

	return nil
}
