package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func EnvEndPoint() string {
	return os.Getenv("ENDPOINT")
}

func EnvAccessKey() string {
	return os.Getenv("ACCESSKEYID")
}

func EnvSecretAccesKey() string {
	return os.Getenv("SECRETACCESSKEY")
}

func EnvBucketName() string {
	return os.Getenv("BUCKETNAME")
}
