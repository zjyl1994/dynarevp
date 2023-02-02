package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigS struct {
	Listen    string
	Prefix    string
	Method    []string
	Whitelist []string
	Blacklist []string
}

var Config ConfigS

func loadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &Config)
}
