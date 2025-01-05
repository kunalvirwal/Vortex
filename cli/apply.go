package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	pb "github.com/kunalvirwal/Vortex/proto/factory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func applyDeployment(applyCmd *flag.FlagSet, f *string) {
	applyCmd.Parse(os.Args[2:])

	if strings.TrimSpace(*f) == "" {
		fmt.Println("Error: File path is required!")
		fmt.Println("Usage: vortex apply -f <path to deployment>")
		os.Exit(1)
	}

	yamlData, err := os.ReadFile(*f)
	if err != nil {
		fmt.Println("Error reading the deployment file:", err)
		os.Exit(1)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to the vortex-service:", err)
		os.Exit(1)
	}
	defer conn.Close()
	client := pb.NewContainerFactoryClient(conn)

	ctx := context.Background()
	res, err := client.Apply(ctx, &pb.RequestBody{
		Data: yamlData,
	})
	if err != nil {
		fmt.Println("Error applying the deployment:", err)
		os.Exit(1)
	}
	if res.GetSuccess() {
		fmt.Println("Deployment applied successfully!")
	} else {
		fmt.Println("Error applying the deployment")
	}

}
