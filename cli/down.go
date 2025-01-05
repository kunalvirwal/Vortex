package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	pb "github.com/kunalvirwal/Vortex/proto/factory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func stopService(downCmd *flag.FlagSet) {

	downCmd.Parse(os.Args[2:])

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to the vortex-service:", err)
		os.Exit(1)
	}
	defer conn.Close()
	client := pb.NewContainerFactoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.Down(ctx, &pb.NameHolder{
		Name: "kill",
	})
	if err != nil {
		fmt.Println("Error stopping the vortex-service:", err)
		os.Exit(1)
	}
	fmt.Println(res.GetName())

}
