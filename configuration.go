package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	StateFile string            `json:"state_file"`
	S3        *S3Config         `json:"s3"`
	Options   map[string]string `json:"options"`
}

type S3Config struct {
	BucketName string `json:"bucket_name"`
	BucketKey  string `json:"bucket_key"`
}

func parseConfiguration(file string) (*Config, error) {
	newConfig := new(Config)
	readFile, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	jsonDecoder := json.NewDecoder(readFile)
	err = jsonDecoder.Decode(newConfig)

	if err != nil {
		return nil, err
	}

	return newConfig, nil
}
