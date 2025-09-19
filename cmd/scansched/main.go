package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/manifoldco/promptui"

	"NetNudge/internal/config"
	"NetNudge/internal/jobs"
	"NetNudge/internal/registry"
)

var asciiBanner = `Net-Nudge`

func main() {
	// CLI flags
	configPath := flag.String("config", "", "Path to configuration YAML file")
	listTools := flag.Bool("list", false, "List all available tools")
	flag.Parse()

	// No arguments? Show menu
	if len(os.Args) == 1 {
		runInteractiveMenu()
		return
	}

	// If --list flag is used
	if *listTools {
		fmt.Println("Available Tools:")
		for _, tool := range registry.ListTools() {
			fmt.Println(" -", tool)
		}
		return
	}

	// Config file required
	if *configPath == "" {
		fmt.Println("Error: --config must be provided")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Register default tools
	registry.RegisterDefaults()

	// Create job scheduler
	scheduler := jobs.NewScheduler()

	// Schedule jobs from config
	for _, job := range cfg.Jobs {
		if err := scheduler.ScheduleJob(job); err != nil {
			log.Printf("Failed to schedule job %s: %v", job.Name, err)
		}
	}

	// Start jobs
	scheduler.Run()
}

// Interactive menu when no args are provided
func runInteractiveMenu() {
	fmt.Println(asciiBanner)

	menu := promptui.Select{
		Label: "Select an option",
		Items: []string{
			"Run scheduled jobs (from config.yaml)",
			"List available tools",
			"Exit",
		},
	}

	_, choice, err := menu.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch choice {
	case "Run scheduled jobs (from config.yaml)":
		// Ask for config file path interactively
		prompt := promptui.Prompt{
			Label:   "Enter config file path",
			Default: "config.yaml",
		}
		configPath, _ := prompt.Run()

		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		registry.RegisterDefaults()
		scheduler := jobs.NewScheduler()

		for _, job := range cfg.Jobs {
			if err := scheduler.ScheduleJob(job); err != nil {
				log.Printf("Failed to schedule job %s: %v", job.Name, err)
			}
		}
		scheduler.Run()

	case "List available tools":
		registry.RegisterDefaults()
		fmt.Println("Available Tools:")
		for _, tool := range registry.ListTools() {
			fmt.Println(" -", tool)
		}

	case "Exit":
		fmt.Println("Goodbye!")
		time.Sleep(time.Second)
		os.Exit(0)
	}
}


