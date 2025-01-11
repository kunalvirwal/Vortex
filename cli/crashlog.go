package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/kunalvirwal/Vortex/proto/factory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func crashlog(crashCmd *flag.FlagSet, d *string, s *string, u *int) {
	crashCmd.Parse(os.Args[2:])
	*d = strings.TrimSpace(*d)
	*s = strings.TrimSpace(*s)
	query := "" // "deploymentName serviceName uid"
	if *d == "" {
		fmt.Println("Error: Deployment name is required!")
		return
	} else if *s == "" {
		fmt.Println("Error: Service name is required!")
		return
	} else if *u < 1 {
		fmt.Println("Error: ServiceUID is Invalid!")
		return
	} else {
		query = *d + " " + *s + " " + fmt.Sprint(*u)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to the vortex-service:", err)
		os.Exit(1)
	}
	defer conn.Close()
	client := pb.NewContainerFactoryClient(conn)

	ctx := context.Background()
	res, err := client.CrashLog(ctx, &pb.NameHolder{
		Name: query,
	})
	if err != nil {
		fmt.Println("Error getting crashlog:", err)
		os.Exit(1)
	}
	data := res.GetData()
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		os.Exit(1)
	}

	fmt.Println("CrashCount:", result["CrashCount"])
	fmt.Println("LastCrashTime:", result["LastCrashTime"])
	fmt.Println("CurrentBackoffDuration:", result["CurrentBackoffDuration"])
	fmt.Println("IsInCrashLoop:", result["IsInCrashLoop"])
	crashHistory, ok := result["CrashHistory"].([]interface{})
	if !ok {
		fmt.Println("Error: CrashHistory is not in the expected format")
		return
	}
	fmt.Println("CrashHistory:")
	for i, val := range crashHistory {
		fmt.Println(strconv.Itoa(i+1)+". Exit Code:", val.(map[string]interface{})["ExitCode"])
		fmt.Println("   Crash Time:", val.(map[string]interface{})["Timestamp"], "sec")
		fmt.Println("   Signal:", val.(map[string]interface{})["Signal"])
		fmt.Println("   Exit Error:", val.(map[string]interface{})["ExitError"])
		fmt.Println("   ----------------------")
	}

}
