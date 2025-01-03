package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
)

func NewClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error creating docker client")
		os.Exit(1)
	}
	_, err = cli.Ping(context.Background())
	if err != nil {
		fmt.Println("Error pinging docker client")
		os.Exit(1)
	}
	fmt.Println("Connection with docker daemon established")
	return cli
}

var cli = NewClient()
