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
