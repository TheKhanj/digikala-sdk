package config

import (
	"encoding/json"
	"os"
)


func ReadConfig(c *Config, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return c.UnmarshalJSON(b)
}

func WriteConfig(c *Config, path string) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0644)
}
