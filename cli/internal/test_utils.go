package internal

import (
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/thekhanj/digikala-sdk/common"
)

func GetAbsPath(name string) string {
	env := os.Getenv("env")

	if !(env == "test" || env == "dev") {
		return name
	}

	if path.IsAbs(name) {
		return name
	}

	return path.Join(common.GetProjectRoot(), name)
}

func AssertAtLeastOneFile(dir string) error {
	entries, err := os.ReadDir(GetAbsPath(dir))
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			return nil
		}
	}

	return fmt.Errorf(
		"no files found in directory: %s", dir,
	)
}

func GetProcDir() string {
	if uid := syscall.Getuid(); uid == 0 {
		return fmt.Sprintf("/var/run/digkala-api/%d", os.Getpid())
	} else {
		return fmt.Sprintf("/run/user/%d/digikala-sdk/%d", uid, os.Getpid())
	}
}

func AssertDir(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return errors.New("path is not a directory: " + path)
		}

		return nil
	}

	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return errors.New("failed to create directory: " + path)
		}
		return nil
	}

	return err
}
