package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/codechamp1/certlens/configs"
	"github.com/codechamp1/certlens/internal/data"
	"github.com/codechamp1/certlens/internal/service"
	"github.com/codechamp1/certlens/internal/ui"
)

func main() {
	config := configs.Load()

	secretsFetcher, err := data.NewDefaultSecretsFetcher(config.KubeConfigPath, config.Context)

	if err != nil {
		log.Fatalf("Failed to create Kubernetes secrets fetcher: %v", err)
	}

	parser := data.NewDefaultParser()

	repo := data.NewDefaultRepository(secretsFetcher, parser)

	manager := service.NewDefaultManager(repo)

	model, err := ui.NewModel(manager, config.Namespace, config.Name)

	if err != nil {
		log.Fatalf("Failed to create UI model: %v", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
