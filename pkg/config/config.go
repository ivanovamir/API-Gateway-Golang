package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"unicode"
)

type config struct {
	Path string
}

type Config interface {
	ParseConfig() (Data, error)
	validate(dst []byte, src *Data) error
	isUpperCase(str string) bool
}

func NewConfig(options ...Option) Config {
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

	if err := c.validate(body, &data); err != nil {
		return Data{}, err
	}

	return data, nil
}

func (c *config) validate(dst []byte, src *Data) error {

	if err := json.Unmarshal(dst, &src); err != nil {
		return fmt.Errorf("error occured validating: %v", err)
	}

	for _, val := range src.Data {
		if val.Path == "" {
			return fmt.Errorf("path is empty")
		}

		if val.Url == "" {
			return fmt.Errorf("url is empty")
		}

		if val.Method == "" {
			return fmt.Errorf("method is empty")
		}

		if !c.isUpperCase(val.Method) {
			return fmt.Errorf("method must be uppercase")
		}

		if val.MakeProxy {
			if val.ProxyUrl == "" {
				return fmt.Errorf("proxy url is empty")
			}
		}
	}

	return nil
}

func (c *config) isUpperCase(str string) bool {
	for _, r := range str {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
