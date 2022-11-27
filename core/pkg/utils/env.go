package utils

import (
	"os"
	"strconv"
)

func RequiredStringEnv(envKey string) string {
	env := os.Getenv(envKey)
	if env == "" {
		panic("Please, set environment variable " + envKey)
	}
	return env
}

func RequiredIntEnv(envKey string) int {
	env := os.Getenv(envKey)
	if env == "" {
		panic("Please, set environment variable " + envKey)
	}
	intEnv, err := strconv.Atoi(env)
	if err != nil {
		panic(err)
	}
	return intEnv
}

func GetFullServiceAddress(serviceName string) string {
	return RequiredStringEnv(serviceName+"_SERVICE_ADDRESS") + ":" + RequiredStringEnv(serviceName+"_SERVICE_PORT")
}
