package config

import "os"

func GetUniversalDir() (string, error) {
	dir, err := os.Getwd()
	return dir, err
}
