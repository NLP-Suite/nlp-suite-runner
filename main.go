package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

const (
	uiImageName    = "hello-world:latest"    // Replace with your actual UI image name
	agentImageName = "hello-world:latest" // Replace with your actual agent image name
	agentSourceFolder = "nlp-suite"
	agentTargetMountPath = "/mnt/lib"
)

func main() {

	os.Setenv("DOCKER_API_VERSION", "1.42")

	ctx := context.Background()

	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		return
	}

	// Pull the latest UI image
	fmt.Println("Installing the latest version of the NLP Suite UI...")
	if err := pullImage(ctx, cli, uiImageName); err != nil {
		fmt.Println("Error pulling UI image:", err)
		return
	}

	// Pull the latest agent image
	fmt.Println("Installing the latest version of the NLP Suite Agent...")
	if err := pullImage(ctx, cli, agentImageName); err != nil {
		fmt.Println("Error pulling agent image:", err)
		return
	}

	fmt.Println("Validating the NLP Suite Folder...")
	targetMountPath, err := validateMountPoint(agentSourceFolder)
	if err != nil {
		fmt.Println("Error creating agent mount point:", err)
		return
	}

	// Run the containers
	fmt.Println("Starting the NLP Suite UI...")
	if err := runContainer(ctx, cli, uiImageName, "", ""); err != nil {
		fmt.Println("Error running UI container:", err)
		return
	}

	fmt.Println("Starting the NLP Suite Agent...")
	if err := runContainer(ctx, cli, agentImageName, targetMountPath, agentTargetMountPath); err != nil {
		fmt.Println("Error running agent container:", err)
		return
	}

	// Wait indefinitely for container exit (or interrupt with Ctrl+C)
	c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
		fmt.Println("Exiting -- Cleaning up the NLP Suite...")
		if err := cleanUp(ctx, cli); err != nil {
			fmt.Println("Error cleaning up:", err)
		}
        os.Exit(0)
    }()

	fmt.Println("The NLP Suite is running... view the UI at http://localhost:8000")
	select {}
}

// Pulls any image
func pullImage(ctx context.Context, cli *client.Client, imageName string) error {
	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.ReadAll(reader)
	return err
}

// Validates that the mount path exists and creates the folder if it does not
func validateMountPoint(sourceFolder string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	mountPath := path.Join(home, sourceFolder)

	if _, err := os.Stat(mountPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(mountPath, 0755); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return mountPath, nil
}

// Runs a container that already exists
func runContainer(ctx context.Context, cli *client.Client, imageName, sourceMountPath string, targetMountPath string) error {
	config := &container.Config{
		Image: imageName,
		// TODO: Add any additional configuration needed for your containers
	}
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				// TODO: Fix
				Source: sourceMountPath,
				Target: targetMountPath,
				Type:   mount.TypeBind,
			},
		},
	}

	if sourceMountPath == "" || targetMountPath == "" {
		hostConfig = nil
	}

	// TODO: Define network configuration if necessary

	c, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, c.ID, container.StartOptions{}); err != nil {
		return err
	}

	return nil
}

// Cleans up running containers
func cleanUp(ctx context.Context, cli *client.Client) error {
	// Get all running containers
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err
	}

	// Stop running containers
	for _, c := range containers {
		if err := cli.ContainerStop(ctx, c.ID, container.StopOptions{}); err != nil {
			fmt.Printf("Error stopping container %s: %v\n", c.ID, err)
		}
	}

	// Remove all containers
	for _, c := range containers {
		if err := cli.ContainerRemove(ctx, c.ID, container.RemoveOptions{Force: true}); err != nil {
			fmt.Printf("Error removing container %s: %v\n", c.ID, err)
		}
	}

	// Remove downloaded images (optional)
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return err
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == uiImageName || tag == agentImageName { // Remove only specific images based on names
				if _, err := cli.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{Force: true}); err != nil {
					fmt.Printf("Error removing image %s: %v\n", image.ID, err)
				}
			}
		}
	}

	return nil
}
