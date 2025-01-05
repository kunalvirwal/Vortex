package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func startService(upCmd *flag.FlagSet) {

	upCmd.Parse(os.Args[2:])

	path := getPath()
	cmd = exec.Command(path)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting the vortex-service:", err)
		os.Exit(1)
	}
	fmt.Println("Started the vortex-service")

}

func getPath() string {
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path")
		return "./vortex-service"
	} else {
		executableDir := filepath.Dir(executablePath)
		path := filepath.Join(executableDir, "vortex-service")
		_, err := os.Stat(path)
		if err != nil {
			fmt.Println("Unable to find vortex-service binary. Rebuild the application")
			os.Exit(1)
		}
		return path
	}
}
