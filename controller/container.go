package controller

import (
	"bytes"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

func (c *Controller) CreateContainerName(code string) string {
	return fmt.Sprintf("bp-container-%s", code)
}

func (c *Controller) CreateNewContainer(ctx context.Context, name, image string) (string, error) {
	containerConfig := &container.Config{
		NetworkDisabled: true,
		Hostname:        name,
		Image:           image,
	}

	hostConfig := &container.HostConfig{
		NetworkMode: "none",
	}

	containerBody, err := c.cli.ContainerCreate(
		ctx, containerConfig, hostConfig, nil, nil, name)
	if err != nil {
		return "", err
	}

	return containerBody.ID, nil
}

func (c *Controller) StartContainerByID(ctx context.Context, id string) error {
	return c.cli.ContainerStart(ctx, id, types.ContainerStartOptions{})
}

func (c *Controller) RemoveContainerByID(ctx context.Context, id string) error {
	return c.cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
}

func (c *Controller) GetContainerStatus(ctx context.Context, containerID string) (string, int, error) {
	container, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", 0, err
	}

	return container.State.Status, container.State.ExitCode, nil
}

func (c *Controller) GetExecutionOutput(ctx context.Context, containerID string) (string, error) {
	executionOutput, err := c.cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	defer executionOutput.Close()
	if err != nil {
		c.l.Errorf("error getting container logs %s: %s", containerID, err)
		return "", err
	}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	stdcopy.StdCopy(stdout, stderr, executionOutput)
	return stdout.String(), nil
}
