package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	profilesManager := NewProfilesManager()
	if err := profilesManager.LoadProfiles(); err != nil {
		log.Printf("Warning: Could not load profiles: %v", err)
	}

	m := NewMainModel(profilesManager)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
