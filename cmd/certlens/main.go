package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"certlens/configs"
	"certlens/internal/service"
	"certlens/internal/ui"
)

func main() {
	config := configs.Load()

	//client, err := client.NewClient(config.KubeConfigPath)
	//
	//if err != nil {
	//	log.Fatalf("Failed to create Kubernetes client: %v", err)
	//}
	//
	//repo := repository.NewSecretsRepository(client)

	//svc := service.NewSecretsService(repo)
	mockSvc := service.NewMockSecretService()

	model, err := ui.NewModel(config.Namespace, mockSvc)

	if err != nil {
		log.Fatalf("Failed to create UI model: %v", err)
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
