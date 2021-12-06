package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	RDSConnectionURL string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	rdsConnectionURL, err := getEnv("DB_CONN_URL")
	if err != nil {
		return nil, err
	}

	return &Config{
		RDSConnectionURL: rdsConnectionURL,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
