package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var cmd *exec.Cmd

func main() {
	// Start the vortex service
	upCmd := flag.NewFlagSet("up", flag.ExitOnError) // TODO : Implement ssh command tio start the service on other servers too

	// Stop the vortex service
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)

	// // Apply a deployment
	applyCmd := flag.NewFlagSet("apply", flag.ExitOnError)
	applyDep := applyCmd.String("f", "", "Path to the deployment to apply")

	// // Delete a deployment
	// deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Usage: vortex <command> <args>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		startService(upCmd)

	case "down":
		stopService(downCmd)

	case "apply":
		applyDeployment(applyCmd, applyDep)
	}

}
