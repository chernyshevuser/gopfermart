package config

import (
	"fmt"
	"os"
)

func GetConfigString(key configKey) (string, error) {
	val := os.Getenv(string(key))
	if len(val) == 0 {
		return "", fmt.Errorf("error in getting %s value", key)
	}

	return val, nil
}
