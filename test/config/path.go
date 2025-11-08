package config

import (
	"log"
	"os"
	"path/filepath"
)

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		if dir == filepath.Dir(dir) {
			log.Fatalf("go.mod not found in any parent directories")
		}

		dir = filepath.Dir(dir)
	}
}

func GetConfigPath() string {
	root := getProjectRoot()
	path, ok := os.LookupEnv("TEST_CONFIG")
	if ok {
		return filepath.Join(root, path)
	}

	return filepath.Join(root, "test-config.json")
}
