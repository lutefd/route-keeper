package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lutefd/route-keeper/internal/models"
	"github.com/lutefd/route-keeper/internal/ui"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
	BuiltBy   = "unknown"
	GoVersion = "unknown"
)

func printVersion() {
	fmt.Printf("Route Keeper - API Monitoring Tool\n")
	fmt.Printf("Version:    %s\n", Version)
	fmt.Printf("Commit:     %s\n", Commit)
	fmt.Printf("Build Date: %s\n", BuildDate)
	fmt.Printf("Built by:   %s\n", BuiltBy)
	fmt.Printf("Go version: %s\n", GoVersion)
	os.Exit(0)
}

func main() {
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *versionFlag {
		printVersion()
	}

	profilesManager := models.NewProfilesManager()
	if err := profilesManager.LoadProfiles(); err != nil {
		log.Printf("Warning: Could not load profiles: %v", err)
	}

	m := ui.NewMainModel(profilesManager)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
