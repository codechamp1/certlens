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

type secretsLoadedMsg struct {
	secrets []list.Item
	err     error
}

type secretLoadedMsg struct {
	secret list.Item
	err    error
}

type secretDelegate struct {
	list.DefaultDelegate
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
	Name          string
	Namespace     string
	SecretService service.SecretsService
	theme         Theme
	width         int
	height        int
}

func NewModel(svc service.SecretsService, namespace, name string) (Model, error) {
	var items []list.Item
	secretsList := list.New(items, newSecretDelegate(), 50, 20)
	secretsList.Title = "Select a TLS Secret"
	return Model{
		Name:          name,
		Namespace:     namespace,
		SecretService: svc,
		secrets:       secretsList,
		theme:         Default,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return loadSecretsCmd(m)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.secrets.FilterState() != list.Filtering {
			switch msg.String() {
			case "u":
				return m, loadSecretsCmd(m)
			}
		}

	case tea.WindowSizeMsg:
		h, v := m.theme.docStyle.GetFrameSize()
		m.width = msg.Width
		m.height = msg.Height
		m.secrets.SetSize(msg.Width-h, msg.Height-v)

	case secretsLoadedMsg:
		if msg.err != nil {
			log.Fatalln(fmt.Errorf("secrets loaded error: %v", msg.err))
		}
		m.secrets.SetItems(msg.secrets)
	case secretLoadedMsg:
		if msg.err != nil {
			log.Fatalln(fmt.Errorf("secrets loaded error: %v", msg.err))
		}
		m.secrets.SetItems([]list.Item{msg.secret})
	}

	var cmd tea.Cmd
	m.secrets, cmd = m.secrets.Update(msg)

	if item, ok := m.secrets.SelectedItem().(secretItem); ok {
		m.selected = &item
	} else {
		m.selected = nil
	}
	return m, cmd
}
func (m Model) View() string {
	h, v := m.theme.docStyle.GetFrameSize()
	usableWidth := m.width - h
	usableHeight := m.height - v

	leftWidth := usableWidth / 2
	rightWidth := usableWidth - leftWidth

	leftView := lipgloss.NewStyle().Width(leftWidth).Height(usableHeight).Render(m.secrets.View())

	rightStyle := lipgloss.NewStyle().Width(rightWidth).Height(usableHeight)
	var rightView string

	if m.selected != nil {
		data, err := m.SecretService.InspectTLSSecret(m.selected.namespace, m.selected.name)
		if err != nil {
			rightView = rightStyle.Render(err.Error())
		} else {
			rightView = rightStyle.Render(formatCertificateInfo(*data, m.theme))
		}
	} else {
		rightContent := fmt.Sprintf("Nothing yet selected, waiting.....")
		rightView = rightStyle.Render(rightContent)
	}

	return m.theme.docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView))
}

func loadSecretsCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		if m.SecretService == nil {
			return secretsLoadedMsg{nil, fmt.Errorf("secrets service not initialized")}
		}

		if m.Name != "" {
			secret, err := m.SecretService.ListTLSSecret(m.Namespace, m.Name)
			return secretLoadedMsg{secretItem{secret.Name, secret.Namespace}, err}
		}

		secrets, err := m.SecretService.ListTLSSecrets(m.Namespace)

		items := make([]list.Item, len(secrets))
		for i, s := range secrets {
			items[i] = secretItem{s.Name, s.Namespace}
		}
		return secretsLoadedMsg{items, err}
	}
}

func newSecretDelegate() secretDelegate {
	delegate := list.NewDefaultDelegate()

	delegate.ShortHelpFunc = func() []key.Binding {
		keys := []key.Binding{key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "refresh secrets"),
		)}
		return keys
	}

	delegate.FullHelpFunc = func() [][]key.Binding {
		groups := [][]key.Binding{{key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "refresh secrets"),
		)}}
		return groups
	}

	// TODO: use the theme styles
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("#00BFFF"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(lipgloss.Color("#5DADE2"))

	return secretDelegate{delegate}
}
