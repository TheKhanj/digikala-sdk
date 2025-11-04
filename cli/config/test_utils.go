package config

import (
	"path"

	"github.com/thekhanj/digikala-sdk/cli/internal"
)

func ReadTestConfig() (*Config, error) {
	return ReadConfig(
		path.Join(internal.GetProjectRoot(), "github-config.json"),
	)
}
