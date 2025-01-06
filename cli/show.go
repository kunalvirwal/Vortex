package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	pb "github.com/kunalvirwal/Vortex/proto/factory"
	"github.com/kunalvirwal/Vortex/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func show(showCmd *flag.FlagSet) {
	showCmd.Parse(os.Args[2:])

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to the vortex-service:", err)
		os.Exit(1)
	}
	defer conn.Close()
	client := pb.NewContainerFactoryClient(conn)

	ctx := context.Background()
	res, err := client.Show(ctx, &pb.NameHolder{
		Name: "all",
	})
	if err != nil {
		fmt.Println("Error stopping the vortex-service:", err)
		os.Exit(1)
	}
	var result types.State
	err = json.Unmarshal(res.GetData(), &result)
	if err != nil {
		fmt.Println("Error unmarshaling the response:", err)
		os.Exit(1)
	}
	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	// Print to terminal
	fmt.Println(string(prettyJSON))
}
