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

func delete(deleteCmd *flag.FlagSet, d *string, s *string) {
	deleteCmd.Parse(os.Args[2:])
	*d = strings.TrimSpace(*d)
	*s = strings.TrimSpace(*s)
	query := "" // "deploymentName" or "deploymentName serviceName"
	if *d != "" {
		var deleting string
		if *s == "" {
			// "Deleting deployment:", *d (all services)
			query = *d
			deleting = "Deployment"
		} else {
			// "Deleting service:", *s from deployment:", *d
			query = *d + " " + *s
			deleting = "Service"
		}

		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("Error connecting to the vortex-service:", err)
			os.Exit(1)
		}
		defer conn.Close()
		client := pb.NewContainerFactoryClient(conn)

		ctx := context.Background()
		res, err := client.Delete(ctx, &pb.NameHolder{
			Name: query,
		})
		if err != nil {
			fmt.Println("Error deleting", deleting, ":", err)
			os.Exit(1)
		}

		if res.GetSuccess() {
			fmt.Println(deleting, "deleted successfully")
		} else {
			fmt.Println("Error deleting", deleting)
		}

	} else {
		fmt.Println("Error: Deployment name is required!")
	}

}
