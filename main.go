package main

import (
	"fmt"
	"os"
)

var Version = "dev"

func main() {
	fmt.Printf("Go Directory Backup version: %s\n", Version)
	fmt.Println("Starting Go Directory Backup...")
	config, err := LoadOrInitConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	fmt.Println("Backup configuration loaded.")
	RunScheduler(config)
}
