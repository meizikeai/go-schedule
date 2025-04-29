package config

import (
	"fmt"
	"os"
	"slices"
)

var env = []string{
	"release",
	"test",
}

func GetMode() string {
	pass := false
	mode := os.Getenv("GO_MODE")

	pass = slices.Contains(env, mode)

	if pass == false {
		mode = "test"
	}

	return mode
}

func IsProduction() bool {
	mode := GetMode()
	result := false

	if mode == "release" {
		result = true
	}

	return result
}

func getKey(k string) string {
	mode := GetMode()
	result := fmt.Sprintf("%s-%s", k, mode)

	return result
}
