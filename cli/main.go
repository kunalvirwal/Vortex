package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var cmd *exec.Cmd

func main() {
	// Start the vortex service
	// upCmd := flag.NewFlagSet("up", flag.ExitOnError) // TODO : Implement ssh command tio start the service on other servers too

	// // Stop the vortex service
	// downCmd := flag.NewFlagSet("down", flag.ExitOnError)

	// // Apply a deployment
	// applyCmd := flag.NewFlagSet("apply", flag.ExitOnError)
	// applyDep := applycmd.String("f", "", "Path to the deployment to apply")

	// // Delete a deployment
	// deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Usage: vortex <command> <args>")
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	switch os.Args[1] {
	case "up":
		startService()
		fmt.Println(cmd)
		go func() {
			<-sig
			fmt.Println("Stopping the vortex-service")
			cmd.Process.Kill()
			os.Exit(0)
		}()
	case "down":
		stopService()
	}

}
