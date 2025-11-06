package common

import "os"

func GetMandatoryEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(key + " environment variable was not passed to the test")
	}
	if value == "" {
		panic(key + " environment variable is empty")
	}

	return value
}
