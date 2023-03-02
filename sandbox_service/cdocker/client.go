package cdocker

import (
	"context"
	"os"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/docker/docker/client"
)

func LoadDockerImageFromTAR(ctx context.Context, imagePath string, dockerClient *client.Client) error {
	dockerImageFile, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer dockerImageFile.Close()
	loadImageResponse, err := dockerClient.ImageLoad(ctx, dockerImageFile, false)
	if err != nil {
		return err
	}
	defer loadImageResponse.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func ClientFromConfig(config *cviper.CustomViper) (
	*client.Client, error,
) {
	sandboxDockerClientVersion := config.GetStringRequired("DOCKER_CODE_RUNNER_CLIENT_VERSION")
	return client.NewClientWithOpts(client.WithVersion(sandboxDockerClientVersion))
}
