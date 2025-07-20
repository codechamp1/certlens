package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"certlens/internal/service"
)

type secretsLoadedMsg struct {
	secrets []list.Item
}

type secretLoadedMsg struct {
	secret list.Item
}

type secretDelegate struct {
	list.DefaultDelegate
}

type loadingStartedMsg struct{}

type errorMsg struct{ err error }

type secretItem struct {
	name      string
	namespace string
}

func (s secretItem) Title() string       { return s.name }
func (s secretItem) Description() string { return "Namespace: " + s.namespace }
func (s secretItem) FilterValue() string { return s.name }

type Model struct {
	loading       bool
	secrets       list.Model
	selected      *secretItem
	spinner       spinner.Model
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
	spinnerModel := spinner.New()
	spinnerModel.Spinner = spinner.Dot
	return Model{
		Name:          name,
		Namespace:     namespace,
		SecretService: svc,
		secrets:       secretsList,
		spinner:       spinnerModel,
		theme:         Default,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, loadSecretsCmd(m))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
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
		m.secrets.SetItems(msg.secrets)
		m.loading = false
	case secretLoadedMsg:
		m.secrets.SetItems([]list.Item{msg.secret})
		m.loading = false
	case loadingStartedMsg:
		m.loading = true
		cmds = append(cmds, m.spinner.Tick)
	case errorMsg:
		m.loading = false
		m.secrets.Title = fmt.Sprintf("Error: %v", msg.err)
	}

	if !m.loading {
		var listCmd tea.Cmd
		m.secrets, listCmd = m.secrets.Update(msg)
		cmds = append(cmds, listCmd)
	}

	if m.loading {
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinCmd)
	}

	updateSelectedItem(&m)

	return m, tea.Batch(cmds...)
}
func (m Model) View() string {
	h, v := m.theme.docStyle.GetFrameSize()
	usableWidth := m.width - h
	usableHeight := m.height - v

	leftWidth := usableWidth / 2
	rightWidth := usableWidth - leftWidth

	left := m.leftPane(leftWidth, usableHeight)
	right := m.rightPane(rightWidth, usableHeight)

	return m.theme.docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, left, right))
}

func loadSecretsCmd(m Model) tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return loadingStartedMsg{} },
		func() tea.Msg {
			if m.SecretService == nil {
				return errorMsg{fmt.Errorf("secrets service not initialized")}
			}

			if m.Name != "" {
				secret, err := m.SecretService.ListTLSSecret(m.Namespace, m.Name)
				if err != nil {
					return errorMsg{fmt.Errorf("failed to load secret %s/%s: %w", m.Namespace, m.Name, err)}
				}
				return secretLoadedMsg{secretItem{secret.Name, secret.Namespace}}
			}

			secrets, err := m.SecretService.ListTLSSecrets(m.Namespace)

			items := make([]list.Item, len(secrets))
			for i, s := range secrets {
				items[i] = secretItem{s.Name, s.Namespace}
			}
			if err != nil {
				return errorMsg{fmt.Errorf("failed to load secrets in namespace %s: %w", m.Namespace, err)}
			}
			return secretsLoadedMsg{items}
		},
	)
}

func newSecretDelegate() secretDelegate {
	delegate := list.NewDefaultDelegate()

	refreshKey := key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "refresh secrets"),
	)

	delegate.ShortHelpFunc = func() []key.Binding { return []key.Binding{refreshKey} }
	delegate.FullHelpFunc = func() [][]key.Binding { return [][]key.Binding{{refreshKey}} }

	// TODO: use the theme styles
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("#00BFFF"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(lipgloss.Color("#5DADE2"))

	return secretDelegate{delegate}
}

func (m Model) leftPane(width, height int) string {
	style := lipgloss.NewStyle().Width(width).Height(height)
	if m.loading {
		return style.Render(m.spinner.View() + " Loading secrets...")
	}
	return style.Render(m.secrets.View())
}

func (m Model) rightPane(width, height int) string {
	style := lipgloss.NewStyle().Width(width).Height(height)
	if m.selected != nil {
		data, err := m.SecretService.InspectTLSSecret(m.selected.namespace, m.selected.name)
		if err != nil {
			return style.Render(err.Error())
		}
		return style.Render(formatCertificateInfo(*data, m.theme))
	}
	return style.Render("Nothing yet selected, waiting.....")
}

func updateSelectedItem(m *Model) {
	if item, ok := m.secrets.SelectedItem().(*secretItem); ok {
		m.selected = item
	} else {
		m.selected = nil
	}
}
