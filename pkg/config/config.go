package config

import (
	"encoding/json"
	"io"
	"os"
)

type config struct {
	Path string
}

func NewConfig(options ...Option) *config {
	cfg := &config{}

	for _, opt := range options {
		opt(cfg)
	}
	return cfg
}

func (c *config) ParseConfig() (Data, error) {
	dataFile, err := os.Open(c.Path)

	if err != nil {
		return Data{}, err
	}

	defer dataFile.Close()

	var data Data

	body, err := io.ReadAll(dataFile)

	if err != nil {
		return Data{}, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return Data{}, err
	}

	return data, nil
}
