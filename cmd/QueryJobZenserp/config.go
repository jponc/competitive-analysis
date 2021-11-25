package main

import (
	"fmt"
	"os"
)

// Config
type Config struct {
	ZenserpApiKey    string
	RDSConnectionURL string
	AWSRegion        string
	SNSPrefix        string
}

// NewConfig initialises a new config
func NewConfig() (*Config, error) {
	rdsConnectionURL, err := getEnv("DB_CONN_URL")
	if err != nil {
		return nil, err
	}

	awsRegion, err := getEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}

	snsPrefix, err := getEnv("SNS_PREFIX")
	if err != nil {
		return nil, err
	}

	zenserpApiKey, err := getEnv("ZENSERP_API_KEY")
	if err != nil {
		return nil, err
	}

	return &Config{
		AWSRegion:        awsRegion,
		SNSPrefix:        snsPrefix,
		RDSConnectionURL: rdsConnectionURL,
		ZenserpApiKey:    zenserpApiKey,
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if v == "" {
		return "", fmt.Errorf("%s environment variable missing", key)
	}

	return v, nil
}
