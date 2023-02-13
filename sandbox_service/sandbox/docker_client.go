package sandbox

import (
	"archive/tar"
	"bytes"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerClient struct {
	client *client.Client
	config *DockerClientConfig
}

type DockerClientConfig struct {
	imageName                   string
	memoryLimitBytes            int64
	networkDisable              bool
	containerNamePrefix         string
	containerGoSourcesDirectory string // /app
}

func NewDockerClientConfig(imageName string, memoryLimitBytes int64, networkDisable bool,
	containerNamePrefix string, containerGoSourcesDirectory string,
) *DockerClientConfig {
	return &DockerClientConfig{
		imageName:                   imageName,
		memoryLimitBytes:            memoryLimitBytes,
		networkDisable:              networkDisable,
		containerNamePrefix:         containerNamePrefix,
		containerGoSourcesDirectory: containerGoSourcesDirectory,
	}
}

func NewDockerClient(client *client.Client, config *DockerClientConfig) Client {
	return &DockerClient{
		client: client,
		config: config,
	}
}

func (c *DockerClient) ExecuteGoCode(ctx context.Context, sourceCode string, userID string) (
	*bytes.Buffer, *bytes.Buffer, error,
) {
	containerName := c.config.containerNamePrefix + userID

	// out, err := c.client.ImagePull(ctx, c.config.imageName, types.ImagePullOptions{})
	// if err != nil {
	// 	return nil, nil, err
	// }
	// defer out.Close()
	// io.Copy(os.Stdout, out)

	exeContainer, err := c.createContainer(ctx, containerName)
	if err != nil {
		return nil, nil, err
	}
	containerID := exeContainer.ID
	defer c.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
	// fmt.Println("Warnings", exeContainer.Warnings) // TODO: log warnings

	reader, err := c.createTARReader(sourceCode)
	if err != nil {
		return nil, nil, err
	}

	err = c.client.CopyToContainer(
		ctx, containerID, c.config.containerGoSourcesDirectory, reader,
		types.CopyToContainerOptions{AllowOverwriteDirWithFile: false},
	)
	if err != nil {
		return nil, nil, err
	}

	if err = c.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return nil, nil, err
	}

	waitResponce, errCh := c.client.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, nil, err
		}
	case <-waitResponce:
	}

	return c.GetProgramOutput(ctx, containerID)
}

func (c *DockerClient) createTARReader(sourceCode string) (
	*bytes.Buffer, error,
) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	hdr := &tar.Header{
		Name: "main.go",
		Mode: 0o644,
		Size: int64(len(sourceCode)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(sourceCode)); err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf.Bytes()), nil
}

func (c *DockerClient) createContainer(ctx context.Context, containerName string) (
	container.CreateResponse, error,
) {
	return c.client.ContainerCreate(
		ctx,
		&container.Config{
			Image:           c.config.imageName,
			AttachStdout:    true,
			AttachStderr:    true,
			NetworkDisabled: c.config.networkDisable,
			// Tty:             false, # ????
			// Cmd:             []string{"ls", "/app"},
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory: c.config.memoryLimitBytes,
				// CPUQuota: int64(50 * 1024 * 1024 * 1024),
			},
		},
		nil, nil,
		containerName,
	)
}

func (c *DockerClient) GetProgramOutput(ctx context.Context, containerID string) (
	/*stdout*/ *bytes.Buffer /*stderr*/, *bytes.Buffer, error,
) {
	logs, err := c.client.ContainerLogs(
		ctx, containerID,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	_, err = stdcopy.StdCopy(stdout, stderr, logs)
	if err != nil {
		return nil, nil, err
	}

	return stdout, stderr, err
}
