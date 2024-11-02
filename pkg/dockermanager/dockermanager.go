package dockermanager

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type ContainerConfig struct {
	Image       string
	Name        string
	Env         []string
	Cmd         []string
	Network     string
	Mounts      []string // For data volumes
	Ports       map[string]string
	MemoryLimit int64
	CPUShares   int64
	AutoRemove  bool
}

func StartContainer(config *ContainerConfig) (string, error) {
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	if _, err := checkImageExists(dockerClient, ctx, config); err != nil {
		out, err := dockerClient.ImagePull(ctx, config.Image, image.PullOptions{})
		if err != nil {
			return "", err
		}
		io.Copy(os.Stdout, out)
	}

	containerConfig := &container.Config{
		Image: config.Image,
		Env:   config.Env,
		Cmd:   config.Cmd,
	}

	mounts := make([]mount.Mount, len(config.Mounts))
	for i, m := range config.Mounts {
		mounts[i] = mount.Mount{
			Type:   mount.TypeVolume,
			Source: "formdata_data",
			Target: m,
		}
	}

	hostConfig := &container.HostConfig{
		// RestartPolicy: container.RestartPolicy{
		// 	Name: container.RestartPolicyAlways,
		// },
		PortBindings: nat.PortMap{},
		Resources: container.Resources{
			Memory:   config.MemoryLimit,
			NanoCPUs: config.CPUShares * 1e9,
		},
		Mounts:     mounts,
		AutoRemove: config.AutoRemove,
	}

	for hostPort, containerPort := range config.Ports {
		portBinding := nat.PortBinding{
			HostIP:   "0.0.0.0",
			HostPort: hostPort,
		}
		hostConfig.PortBindings[nat.Port(containerPort)] = []nat.PortBinding{portBinding}
	}

	networkingConfig := &network.NetworkingConfig{}

	if config.Network != "" {
		networkingConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				config.Network: {},
			},
		}
	}

	resp, err := dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, nil, config.Name)
	if err != nil {
		return "", err
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func StopContainer(containerID string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	timeout := 0 // Immediate stop
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return err
	}

	if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{}); err != nil {
		return err
	}

	return nil
}

func checkImageExists(client *client.Client, ctx context.Context, config *ContainerConfig) (string, error) {
	images, err := client.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return "", err
	}

	imageExists := false
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == config.Image {
				imageExists = true
				break
			}
		}
		if imageExists {
			break
		}
	}

	if !imageExists {
		return "", fmt.Errorf("image %s not found locally", config.Image)
	}

	return config.Image, nil

}
