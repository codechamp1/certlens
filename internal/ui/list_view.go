package ui

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/codechamp1/certlens/internal/service"
)

type Pane int

const (
	LeftPane Pane = iota
	RightPane
)

type secretsLoadedMsg struct {
	secrets []list.Item
}

type secretLoadedMsg struct {
	secret list.Item
}

type inspectTLSSecretMsg struct{}
type loadingStartedMsg struct{}

type loadSecretsMsg struct{}

type copyMsg struct {
	key bool
}

type switchCertViewMsg struct{}

type switchPaneMsg struct{}

type errorMsg struct{ err error }

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
	certViews         []string
	certPages         paginator.Model
	errorModalMsg     string
	inspectedViewport viewport.Model
	layout            modelLayout
	loading           bool
	inspectRaw        bool
	selectedPane      Pane
	secrets           list.Model
	selected          *secretItem
	spinner           spinner.Model
	inspectedError    error
	name              string
	namespace         string
	secretService     service.SecretsService
	theme             ThemeProvider
}

func NewModel(svc service.SecretsService, namespace, name string) (Model, error) {
	var items []list.Item
	secretsList := list.New(items, newSecretDelegate(), 50, 20)
	secretsList.Title = "Select a TLS Secret"
	return Model{
		certPages:         paginator.New(),
		inspectedViewport: viewport.New(50, 20), // Will be updated later,
		name:              name,
		namespace:         namespace,
		secretService:     svc,
		secrets:           secretsList,
		selectedPane:      LeftPane,
		spinner:           spinner.New(),
		theme:             Default,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg { return loadSecretsMsg{} })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.errorModalMsg != "" {
			return m, tea.Quit
		}

		keyStr := msg.String()

		if keyStr == "ctrl+c" {
			return m, tea.Quit
		}
		if m.selectedPane == RightPane {
			switch keyStr {
			case "left":
				m.certPages.PrevPage()
				m.inspectedViewport.SetContent(m.certViews[m.certPages.Page] + "\n\n" + m.certPages.View())
			case "right":
				m.certPages.NextPage()
				m.inspectedViewport.SetContent(m.certViews[m.certPages.Page] + "\n\n" + m.certPages.View())
			}
		}
		if m.secrets.FilterState() != list.Filtering {
			switch keyStr {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "u":
				cmds = append(cmds, func() tea.Msg { return loadSecretsMsg{} })
			case "tab", "p":
				cmds = append(cmds, func() tea.Msg { return switchPaneMsg{} })
			case "r":
				cmds = append(cmds, func() tea.Msg { return switchCertViewMsg{} })
			case "c":
				cmds = append(cmds, func() tea.Msg { return copyMsg{} })
			case "C":
				cmds = append(cmds, func() tea.Msg { return copyMsg{key: true} })
			}
		}

	case tea.WindowSizeMsg:
		m.updateLayout(msg.Width, msg.Height)
	case copyMsg:
		var copyData string
		tlsCert, tlsKey, err := m.secretService.RawInspectTLSSecret(m.selected.namespace, m.selected.name)
		if err != nil {
			m.errorModalMsg = fmt.Sprintf("Error copying secret: %v", err)
		}
		switch {
		case msg.key:
			copyData = tlsKey
		default:
			copyData = tlsCert
		}
		if err := clipboard.WriteAll(copyData); err != nil {
			m.errorModalMsg = fmt.Sprintf("Error copying secret: %v", err)
		}
	case secretsLoadedMsg:
		m.secrets.SetItems(msg.secrets)
		m.loading = false
	case secretLoadedMsg:
		m.secrets.SetItems([]list.Item{msg.secret})
		m.loading = false
	case switchCertViewMsg:
		m.inspectRaw = !m.inspectRaw
		cmds = append(cmds, func() tea.Msg { return inspectTLSSecretMsg{} })
	case switchPaneMsg:
		m.selectedPane = nextPane(m.selectedPane)
	case loadSecretsMsg:
		cmds = append(cmds, loadSecretsCmd(m))
	case loadingStartedMsg:
		m.loading = true
		cmds = append(cmds, m.spinner.Tick)
	case inspectTLSSecretMsg:
		data, err := m.inspectedTLSSecretContent(m.selected.namespace, m.selected.name, m.inspectRaw)
		m.certViews = data
		m.inspectedError = err
		m.certPages.SetTotalPages(len(data))
		m.certPages.Page = 0
		m.inspectedViewport.SetContent(m.certViews[m.certPages.Page] + "\n\n" + m.certPages.View())
	case errorMsg:
		m.loading = false
		m.errorModalMsg = fmt.Sprintf("Error: %v", msg.err)
	}

	if m.loading {
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinCmd)
	} else {
		switch m.selectedPane {
		case LeftPane:
			var listCmd tea.Cmd
			m.secrets, listCmd = m.secrets.Update(msg)
			cmds = append(cmds, listCmd)

		case RightPane:
			var vpCmd tea.Cmd
			m.inspectedViewport, vpCmd = m.inspectedViewport.Update(msg)
			cmds = append(cmds, vpCmd)
		}
	}

	if sel := m.secrets.SelectedItem(); sel != nil {
		if item, ok := sel.(secretItem); ok {
			if m.selected == nil || item.name != m.selected.name || item.namespace != m.selected.namespace {
				m.selected = &item
				cmds = append(cmds, func() tea.Msg { return inspectTLSSecretMsg{} })
			}
		}
	}

	return m, tea.Batch(cmds...)
}
func (m Model) View() string {
	left := m.leftPane(m.layout.LeftPaneWidth, m.layout.UsableHeight)
	right := m.rightPane(m.layout.RightPaneWidth, m.layout.UsableHeight)

	if m.errorModalMsg != "" {
		return m.renderErrorModal(m.errorModalMsg)
	}

	return m.theme.DocStyle().Render(lipgloss.JoinHorizontal(lipgloss.Top, left, right))
}

func loadSecretsCmd(m Model) tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return loadingStartedMsg{} },
		func() tea.Msg {
			if m.secretService == nil {
				return errorMsg{fmt.Errorf("secrets service not initialized")}
			}

			if m.name != "" {
				secret, err := m.secretService.ListTLSSecret(m.namespace, m.name)
				if err != nil {
					return errorMsg{fmt.Errorf("failed to load secret %s/%s: %w", m.namespace, m.name, err)}
				}
				return secretLoadedMsg{secretItem{secret.Name, secret.Namespace}}
			}

			secrets, err := m.secretService.ListTLSSecrets(m.namespace)
			if err != nil {
				return errorMsg{fmt.Errorf("failed to load secrets in namespace %s: %w", m.namespace, err)}
			}

			items := make([]list.Item, len(secrets))
			for i, s := range secrets {
				items[i] = secretItem{s.Name, s.Namespace}
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

	switchCertViewKey := key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "switch between raw and formatted view"),
	)

	copyCertKey := key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy certificate to clipboard"),
	)

	delegate.ShortHelpFunc = func() []key.Binding { return []key.Binding{refreshKey} }
	delegate.FullHelpFunc = func() [][]key.Binding { return [][]key.Binding{{refreshKey, switchCertViewKey, copyCertKey}} }

	// TODO: use the theme styles
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("#00BFFF"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(lipgloss.Color("#5DADE2"))

	return secretDelegate{delegate}
}

func (m Model) leftPane(width, height int) string {
	style := m.theme.Pane(m.selectedPane == LeftPane, width, height)
	if m.loading {
		return style.Render(m.spinner.View() + " Loading secrets...")
	}
	return style.Render(m.secrets.View())
}

func (m Model) rightPane(width, height int) string {
	style := m.theme.Pane(m.selectedPane == RightPane, width, height)
	if m.inspectedError != nil {
		return style.Render(fmt.Errorf("error inspecting secret: %w", m.inspectedError).Error())
	}
	if m.selected != nil && !m.loading {
		return style.Render(m.inspectedViewport.View())
	}
	return style.Render("Nothing yet selected, waiting.....")
}

func (m Model) renderErrorModal(msg string) string {
	errorModalRender := m.theme.ErrorModalWithWidth(m.layout.TotalWidth).Render("Error:\n" + msg + "\n\nPress any key to dismiss.")
	return lipgloss.Place(m.layout.TotalWidth, m.layout.TotalHeight, lipgloss.Center, lipgloss.Center, errorModalRender)
}

func nextPane(currentPane Pane) Pane {
	if currentPane == LeftPane {
		return RightPane
	}
	return LeftPane
}

func (m Model) inspectedTLSSecretContent(namespace, name string, raw bool) ([]string, error) {
	if raw {
		tlsCert, tlsKey, err := m.secretService.RawInspectTLSSecret(namespace, name)
		if err != nil {
			return nil, fmt.Errorf("failed to inspect secret %s/%s: %w", namespace, name, err)
		}
		return []string{tlsCert, tlsKey}, nil // o singură pagină
	}

	certs, err := m.secretService.InspectTLSSecret(namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect secret %s/%s: %w", namespace, name, err)
	}

	var views []string
	for _, cert := range certs {
		views = append(views, formatCertificateInfo(cert, m.theme))
	}
	return views, nil
}
