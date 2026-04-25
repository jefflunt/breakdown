package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"breakdown/pkg/config"
	"breakdown/pkg/llm"
	"breakdown/pkg/breakdown"
	"breakdown/pkg/version"
)

func isGitRepo() bool {
	cmd := exec.Command("git", "status")
	err := cmd.Run()
	return err == nil
}

func main() {
	if !isGitRepo() {
		fmt.Println("Error: breakdown must be run from inside a Git repository. This ensures accurate codebase context via .gitignore.")
		os.Exit(1)
	}

	// Default configuration
	configFile := config.DefaultPath()

	var planName string
	flag.StringVar(&planName, "plan", "", "Name of the plan to load or create")
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.Parse()

	// Handle sub-commands
	if flag.NArg() > 0 {
		switch flag.Arg(0) {
		case "version":
			fmt.Println(version.Version)
			os.Exit(0)
		case "help":
			flag.Usage()
			os.Exit(0)
		}
	}

	// Task execution logic
	var initialTask string
	// Check if data is piped to STDIN
	stat, err := os.Stdin.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		// Read from STDIN
		data, err := io.ReadAll(os.Stdin)
		if err == nil {
			initialTask = strings.TrimSpace(string(data))
		}
	} else if flag.NArg() > 0 {
		// If it's not a subcommand, treat it as the task
		initialTask = strings.Join(flag.Args(), " ")
	}

	if initialTask == "" {
		fmt.Println("Error: No task provided.")
		os.Exit(1)
	}

	// Load or create default configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Create an overarching context
	ctx := context.Background()

	// Instantiate the LLM based on the loaded config
	client, err := llm.NewClient(ctx, cfg)
	if err != nil {
		fmt.Printf("Failed to initialize LLM client: %v\n", err)
		os.Exit(1)
	}

	// Ensure plans directory exists
	if err := os.MkdirAll(cfg.PlansDir, 0755); err != nil {
		fmt.Printf("Error creating plans directory: %v\n", err)
		os.Exit(1)
	}

	stateFile := filepath.Join(cfg.PlansDir, "state.json")

	p := breakdown.NewPlanner(breakdown.Config{
		PlansDir:  cfg.PlansDir,
		StateFile: stateFile,
		Workspace: "./workspace",
		Verbose:   verbose,
	}, client)

	if err := p.Start(ctx, initialTask); err != nil {
		fmt.Printf("Error during planning: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generating file structure...")
	if err := p.Root.GenerateFilesystemStructure("./breakdown-output"); err != nil {
		fmt.Printf("Error generating file structure: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done. Breakdown output available in ./breakdown-output.")
}
