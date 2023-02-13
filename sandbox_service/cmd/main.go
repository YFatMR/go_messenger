package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/execontroller"
	"github.com/YFatMR/go_messenger/sandbox_service/exeservice"
	"github.com/YFatMR/go_messenger/sandbox_service/grpcserver"
	"github.com/YFatMR/go_messenger/sandbox_service/sandbox"
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := cviper.New()
	config.AutomaticEnv()

	programExecutionTimeout := config.GetMillisecondsDurationRequired("SANDBOX_PROGRAM_EXECUTION_TIMOUT_MILLISECONDS")
	serviceAddress := config.GetStringRequired("SERVICE_ADDRESS")
	sandboxDockerImage := config.GetStringRequired("SANDBOX_DOCKER_IMAGE_NAME")
	sandboxMemoryLimit := config.GetInt64Required("SANDBOX_MEMORY_LIMITATION_BYTES")
	sandboxNetworkDisabled := config.GetBoolRequired("SANDBOX_NETWORK_DISABLED")
	sandboxContainerNamePrefix := config.GetStringRequired("SANDBOX_CONTAINER_NAME_PREFIX")
	sandboxImageSourceDirectory := config.GetStringRequired("SANDBOX_GO_IMAGE_SOURCE_DIRECTORY")
	sandboxDockerClientVersion := config.GetStringRequired("SANDBOX_DOCKER_CLIENT_VERSION")
	sandboxImagePath := config.GetStringRequired("SANDBOX_DOCKER_IMAGE_PATH")

	dockerConfig := sandbox.NewDockerClientConfig(
		sandboxDockerImage, sandboxMemoryLimit, sandboxNetworkDisabled,
		sandboxContainerNamePrefix, sandboxImageSourceDirectory,
	)

	dockerClient, err := client.NewClientWithOpts(client.WithVersion(sandboxDockerClientVersion))
	if err != nil {
		panic(err)
	}
	defer dockerClient.Close()

	// Unpacking docker image
	dockerImageFile, err := os.Open(sandboxImagePath)
	if err != nil {
		panic(err)
	}
	defer dockerImageFile.Close()
	loadImageResponse, err := dockerClient.ImageLoad(ctx, dockerImageFile, false)
	if err != nil {
		panic(err)
	}
	defer loadImageResponse.Body.Close()
	if err != nil {
		panic(err)
	}

	// We check that the desired docker image has been found
	_, _, err = dockerClient.ImageInspectWithRaw(ctx, sandboxDockerImage)
	if err != nil {
		panic(err)
	}

	sandboxClient := sandbox.NewDockerClient(dockerClient, dockerConfig)
	service := exeservice.New(sandboxClient, programExecutionTimeout)
	controller := execontroller.New(service)
	sandboxServer := grpcserver.New(controller)

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		panic(err)
	}
	proto.RegisterSandboxServer(grpcServer, &sandboxServer)
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
