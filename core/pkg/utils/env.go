package utils

import "os"

func GetEnv(envKey string) string {
	env := os.Getenv(envKey)
	if env == "" {
		panic("Please, set environment variable " + envKey)
	}
	return env
}

func GetFullServiceAddress(serviceName string) string {
	return GetEnv(serviceName+"_SERVICE_ADDRESS") + ":" + GetEnv(serviceName+"_SERVICE_PORT")
}
