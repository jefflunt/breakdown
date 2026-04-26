package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	 "github.com/jefflunt/breakdown/pkg/config"
	 "github.com/jefflunt/breakdown/pkg/llm"
	 "github.com/jefflunt/breakdown/pkg/breakdown"
	"github.com/jefflunt/breakdown/pkg/version"
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
	
	if flag.NArg() < 2 {
		fmt.Println("Error: Missing arguments. Usage: breakdown <path-to-prompt-file> <output-folder>")
		os.Exit(1)
	}

	promptFilePath := flag.Arg(0)
	outputFolder := flag.Arg(1)

	data, err := os.ReadFile(promptFilePath)
	if err != nil {
		fmt.Printf("Error reading prompt file '%s': %v\n", promptFilePath, err)
		os.Exit(1)
	}
	
	initialTask = strings.TrimSpace(string(data))

	if initialTask == "" {
		fmt.Printf("Error: Prompt file '%s' is empty.\n", promptFilePath)
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

	// Ensure output directory exists or handle existing directory
	info, err := os.Stat(outputFolder)
	if err == nil {
		if !info.IsDir() {
			fmt.Printf("Error: Path '%s' exists and is not a directory.\n", outputFolder)
			os.Exit(1)
		}
		
		fmt.Printf("Output directory '%s' already exists. Overwrite? [y/N]: ", outputFolder)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Operation cancelled.")
			os.Exit(0)
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(outputFolder, 0755); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Error checking output directory: %v\n", err)
		os.Exit(1)
	}

	p := breakdown.NewPlanner(breakdown.Config{
		Workspace: "./workspace",
		Verbose:   verbose,
	}, client)

	if err := p.Start(ctx, initialTask); err != nil {
		fmt.Printf("Error during planning: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generating file structure...")
	if err := p.Root.GenerateFilesystemStructure(outputFolder); err != nil {
		fmt.Printf("Error generating file structure: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Done. Breakdown output available in %s.\n", outputFolder)
}
