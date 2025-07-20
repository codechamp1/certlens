package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"certlens/configs"
	"certlens/internal/client"
	"certlens/internal/repository"
	"certlens/internal/service"
	"certlens/internal/ui"
)

func main() {
	config := configs.Load()

	kubeClient, err := client.NewClient(config.KubeConfigPath, config.Context)

	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	repo := repository.NewSecretsRepository(kubeClient)

	svc := service.NewSecretsService(repo)

	model, err := ui.NewModel(svc, config.Namespace, config.Name)

	if err != nil {
		log.Fatalf("Failed to create UI model: %v", err)
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
