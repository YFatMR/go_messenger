package main

import (
	"context"
	"flag"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var dockerComposeEnvFile string
	flag.StringVar(&dockerComposeEnvFile, "docker-compose-env-file", "", "")

	flag.Parse()

	command := exec.CommandContext(ctx, "docker-compose", "--env-file", dockerComposeEnvFile, "up")
	// command.Start()
	err := command.Run()
	if err != nil {
		panic(err)
	}
}
