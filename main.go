package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	uiImageName              = "ghcr.io/nlp-suite/nlp-suite-ui:main"
	uiIpAddress              = "172.16.0.10"
	uiPort                   = "8000"
	agentImageName           = "ghcr.io/nlp-suite/nlp-suite-agent:main"
	agentIpAddress           = "172.16.0.11"
	agentPort                = "3000"
	stanfordCoreNlpImageName = "ghcr.io/nlp-suite/stanford-corenlp-docker:master"
	stanfordCoreNlpIpAddress = "172.16.0.12"
	stanfordCoreNlpPort      = "9000"
	agentSourceFolder        = "nlp-suite"
	agentTargetMountPath     = "/root/nlp-suite"
	networkName              = "nlp-suite-network"
	subnet                   = "172.16.0.0/16"
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

	defer func() { cleanUp(ctx, cli, true) }()

	if os.Getenv("ENV") != "dev" {
		// Pull the latest UI image
		fmt.Println("Installing the latest version of the NLP Suite UI...")
		if err := pullImage(ctx, cli, uiImageName); err != nil {
			fmt.Println("Error pulling UI image:", err)
			return
		}

		fmt.Println("Cleaning up any previous artifacts of the NLP Suite...")
		cleanUp(ctx, cli, false)

		// Pull the latest agent image
		fmt.Println("Installing the latest version of the NLP Suite Agent...")
		if err := pullImage(ctx, cli, agentImageName); err != nil {
			fmt.Println("Error pulling agent image:", err)
			return
		}

		// Pull the latest stanford core nlp image
		fmt.Println("Installing the latest version of Stanford CoreNLP...")
		if err := pullImage(ctx, cli, stanfordCoreNlpImageName); err != nil {
			fmt.Println("Error pulling stanford corenlp image:", err)
			return
		}
	} else {
		fmt.Println("Skipping image installation because detected `dev` flag...")
	}

	fmt.Println("Validating the NLP Suite Folder...")
	targetMountPath, err := validateMountPoint(agentSourceFolder)
	if err != nil {
		fmt.Println("Error creating agent mount point:", err)
		return
	}

	fmt.Println("Validating the NLP Suite Input Folder...")
	_, err = validateMountPoint(path.Join(agentSourceFolder, "input"))
	if err != nil {
		fmt.Println("Error creating agent input point:", err)
		return
	}

	fmt.Println("Validating the NLP Suite Output Folder...")
	_, err = validateMountPoint(path.Join(agentSourceFolder, "output"))
	if err != nil {
		fmt.Println("Error creating agent output point:", err)
		return
	}

	fmt.Println("Validating the NLP Suite CSV Input Folder...")
	_, err = validateMountPoint(path.Join(agentSourceFolder, "csvInput"))
	if err != nil {
		fmt.Println("Error creating csvInput folder:", err)
		return
	}

	// Create the network
	fmt.Println("Creating the NLP Suite network...")
	if err := createNetwork(ctx, cli); err != nil {
		fmt.Println("Error creating network:", err)
		return
	}

	// Run the containers
	fmt.Println("Starting the NLP Suite UI...")
	if err := runContainer(ctx, cli, uiImageName, "", "", uiIpAddress, uiPort); err != nil {
		fmt.Println("Error running UI container:", err)
		return
	}

	fmt.Println("Starting the Stanford CoreNLP Server...")
	if err := runContainer(ctx, cli, stanfordCoreNlpImageName, "", "", stanfordCoreNlpIpAddress, stanfordCoreNlpPort); err != nil {
		fmt.Println("Error running stanford corenlp server container:", err, ". The NLP Suite is continuing execution as some tools can be used without it.")
	}

	fmt.Println("Starting the NLP Suite Agent...")
	if err := runContainer(ctx, cli, agentImageName, targetMountPath, agentTargetMountPath, agentIpAddress, agentPort); err != nil {
		fmt.Println("Error running agent container:", err)
		return
	}

	// Wait indefinitely for container exit (or interrupt with Ctrl+C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanUp(ctx, cli, true)
		os.Exit(0)
	}()

	if home, err := os.UserHomeDir(); err != nil {
		fmt.Println("Could not find the location of your home folder!")
		return
	} else {
		fmt.Printf("The NLP Suite is running... copy the following address to a browser to open the NLP Suite at http://127.0.0.1:8000\nYour NLP Suite folder can be found at: %s\n", path.Join(home, agentSourceFolder))
	}
	select {}
}

// Pulls any image
func pullImage(ctx context.Context, cli *client.Client, imageName string) error {
	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)
	return err
}

// Validates that the mount path exists and creates the folder if it does not
func validateMountPoint(folder string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	mountPath := path.Join(home, folder)

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
func runContainer(ctx context.Context, cli *client.Client, imageName, sourceMountPath, targetMountPath, ip, port string) error {
	config := &container.Config{
		Image: imageName,
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Source: sourceMountPath,
				Target: targetMountPath,
				Type:   mount.TypeBind,
			},
		},
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%s/tcp", port)): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: port}},
		},
		NetworkMode: container.NetworkMode(networkName),
	}

	netConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {
				IPAMConfig: &network.EndpointIPAMConfig{
					IPv4Address: ip,
				},
			},
		},
	}

	if sourceMountPath == "" || targetMountPath == "" {
		hostConfig.Mounts = nil
	}

	containerTokens := strings.Split(imageName, "/")
	containerName := strings.Trim(strings.ReplaceAll(strings.Split(containerTokens[len(containerTokens)-1], ":")[0], "-", "_"), "/")

	c, err := cli.ContainerCreate(ctx, config, hostConfig, netConfig, nil, containerName)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, c.ID, container.StartOptions{}); err != nil {
		return err
	}

	return nil
}

// Cleans up running containers
func cleanUp(ctx context.Context, cli *client.Client, isExit bool) error {
	fmt.Println("Exiting -- Cleaning up the NLP Suite...")
	// Get all running containers
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		fmt.Println("Error cleaning up:", err)
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

	cli.NetworkRemove(ctx, networkName)

	if isExit {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("The NLP Suite has successfully closed... please type ENTER to close this window.")
		reader.ReadString('\n')
	}

	return nil
}

// Create the network
func createNetwork(ctx context.Context, cli *client.Client) error {
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return err
	}

	needsNetwork := true
	for _, network := range networks {
		if network.Name == networkName {
			needsNetwork = false
		}
	}

	if needsNetwork {
		fmt.Println("Creating NLP Suite Network")
		_, err := cli.NetworkCreate(ctx, networkName, types.NetworkCreate{
			Driver: "bridge",
			IPAM: &network.IPAM{
				Config: []network.IPAMConfig{
					{
						Subnet: subnet,
					},
				},
			},
			Attachable: true,
		})
		return err
	}

	fmt.Println("Skipping NLP Suite Network... Already created")
	return nil
}
