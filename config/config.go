package config

import (
	"fmt"
	"os"
)

func getMode() string {
	mode := os.Getenv("GO_MODE")

	if mode == "" {
		mode = "test"
	}

	return mode
}

func isProduction() bool {
	mode := getMode()
	result := false

	if mode == "release" {
		result = true
	}

	return result
}

func getKey(k string) string {
	mode := getMode()
	result := fmt.Sprintf("%s-%s", k, mode)

	return result
}
