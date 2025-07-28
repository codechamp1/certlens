package ui

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/codechamp1/certlens/internal/domains"
	"github.com/codechamp1/certlens/internal/service"
)

type Pane int

const (
	LeftPane Pane = iota
	RightPane
)

type inspectTLSSecretMsg struct {
	name      string
	namespace string
}

type loadingStartedMsg struct{}

type loadSecretsMsg struct{}

type copyMsg struct {
	key bool
}

type switchPaneMsg struct{}

type errorMsg struct{ err error }

type secretDelegate struct {
	list.DefaultDelegate
}

type Model struct {
	errorModalMsg string
	layout        modelLayout
	loading       bool
	selectedPane  Pane
	selected      *secretItem
	spinner       spinner.Model
	name          string
	namespace     string
	secretService service.SecretsService
	theme         ThemeProvider
	listView      ListViewModel
	certView      CertViewModel
}

func NewModel(svc service.SecretsService, namespace, name string) (Model, error) {
	return Model{
		name:          name,
		namespace:     namespace,
		secretService: svc,
		selectedPane:  LeftPane,
		spinner:       spinner.New(),
		theme:         Default,
		listView:      NewListViewModel(),
		certView:      NewCertViewModel(svc, Default),
	}, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg { return loadSecretsMsg{} })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var listCmds []tea.Cmd
	var certCmds []tea.Cmd
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

		if m.listView.secrets.FilterState() != list.Filtering {
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
		m.loading = false
	case switchPaneMsg:
		m.selectedPane = nextPane(m.selectedPane)
	case loadSecretsMsg:
		cmds = append(cmds, loadSecretsCmd(m))
	case loadingStartedMsg:
		m.loading = true
		cmds = append(cmds, m.spinner.Tick)
	case errorMsg:
		m.loading = false
		m.errorModalMsg = fmt.Sprintf("Error: %v", msg.err)
	case selectedSecretMsg:
		if msg.name != m.selected.name || msg.namespace != m.selected.namespace {
			m.selected = &secretItem{name: msg.name, namespace: msg.namespace}
			certCmds = append(certCmds, func() tea.Msg { return inspectTLSSecretMsg{msg.namespace, msg.name} })
		}
	}

	if m.loading {
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinCmd)
	} else {
		switch m.selectedPane {
		case LeftPane:
			var listCmd tea.Cmd
			m.listView, listCmd = m.listView.Update(msg)
			cmds = append(cmds, listCmd)
		case RightPane:
			var certCmd tea.Cmd
			m.certView, certCmd = m.certView.Update(msg)
			cmds = append(cmds, certCmd)
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
				return updateSecretsMsg{[]domains.K8SResourceID{secret}}
			}

			secrets, err := m.secretService.ListTLSSecrets(m.namespace)
			if err != nil {
				return errorMsg{fmt.Errorf("failed to load secrets in namespace %s: %w", m.namespace, err)}
			}

			return updateSecretsMsg{secrets}
		},
	)
}

func (m Model) leftPane(width, height int) string {
	style := m.theme.Pane(m.selectedPane == LeftPane, width, height)
	if m.loading {
		return style.Render(m.spinner.View() + " Loading secrets...")
	}
	return style.Render(m.listView.View())
}

func (m Model) rightPane(width, height int) string {
	style := m.theme.Pane(m.selectedPane == RightPane, width, height)
	if m.selected != nil && !m.loading {
		return style.Render(m.certView.View())
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
