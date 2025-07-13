package ui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"certlens/internal/service"
)

var docStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1)

type secretsLoadedMsg struct {
	secrets []list.Item
	err     error
}

type secretDelegate struct {
	list.DefaultDelegate
}

func newSecretDelegate() secretDelegate {
	return secretDelegate{
		DefaultDelegate: list.NewDefaultDelegate(),
	}
}

type secretItem struct {
	name      string
	namespace string
}

func (s secretItem) Title() string       { return s.name }
func (s secretItem) Description() string { return "Namespace: " + s.namespace }
func (s secretItem) FilterValue() string { return s.name }

type Model struct {
	secrets       list.Model
	selected      *secretItem
	Namespace     string
	SecretService service.SecretsService
	width         int
	height        int
}

func NewModel(namespace string, svc service.SecretsService) (Model, error) {
	var items []list.Item
	secretsList := list.New(items, newSecretDelegate(), 50, 20)
	secretsList.Title = "Select a TLS Secret"
	return Model{
		Namespace:     namespace,
		SecretService: svc,
		secrets:       secretsList,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return loadSecretsCmd(m)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "u":
			return m, loadSecretsCmd(m)
		case "enter":
			if item, ok := m.secrets.SelectedItem().(secretItem); ok {
				m.selected = &item
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.width = msg.Width
		m.height = msg.Height
		m.secrets.SetSize(msg.Width-h, msg.Height-v)

	case secretsLoadedMsg:
		if msg.err != nil {
			log.Fatalln(fmt.Errorf("secrets loaded error: %v", msg.err))
		}
		m.secrets.SetItems(msg.secrets)
	}

	var cmd tea.Cmd
	m.secrets, cmd = m.secrets.Update(msg)
	return m, cmd
}
func (m Model) View() string {
	h, v := docStyle.GetFrameSize()
	usableWidth := m.width - h
	usableHeight := m.height - v

	leftWidth := usableWidth / 2
	rightWidth := usableWidth - leftWidth

	leftView := lipgloss.NewStyle().Width(leftWidth).Height(usableHeight).Render(m.secrets.View())

	rightStyle := lipgloss.NewStyle().Width(rightWidth).Height(usableHeight)
	var rightView string
	if m.selected != nil {
		rightContent := fmt.Sprintf("Name: %s\nNamespace: %s\n", m.selected.name, m.selected.namespace)
		rightView = rightStyle.Render(rightContent)
	} else {
		rightContent := fmt.Sprintf("Nothing yet selected, waiting.....")
		rightView = rightStyle.Render(rightContent)
	}

	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView))
}

func loadSecretsCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		if m.SecretService == nil {
			return secretsLoadedMsg{nil, fmt.Errorf("secrets service not initialized")}
		}
		secrets, err := m.SecretService.ListTLSSecrets(m.Namespace)
		items := make([]list.Item, len(secrets))
		for i, s := range secrets {
			items[i] = secretItem{s.Name, s.Namespace}
		}
		return secretsLoadedMsg{items, err}
	}
}

func (sd secretDelegate) AdditionalShortHelpKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithHelp("u", "refresh secrets"),
		),
	}
}
