package ui

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/codechamp1/certlens/internal/domains/tls"
	"github.com/codechamp1/certlens/internal/service"
)

type Pane int

const (
	LeftPane Pane = iota
	RightPane
)

type secretsLoadedMsg struct {
	secrets []tls.Secret
}

type renderSecretsListMsg struct{}

type inspectTLSSecretMsg struct {
	tag int
}

type startLoadingMsg struct{}

type loadSecretsMsg struct{}

type copyMsg struct {
	key bool
}

type switchCertViewMsg struct{}

type switchPaneMsg struct{}

type errorMsg struct{ err error }

type secretsListDelegate struct {
	list.DefaultDelegate
}

type secretsListItem struct {
	tls.Secret
}

func (s secretsListItem) Title() string       { return s.Name() }
func (s secretsListItem) Description() string { return "Namespace: " + s.Namespace() }
func (s secretsListItem) FilterValue() string { return s.Name() }

const debounceDuration = 100 * time.Millisecond

type Model struct {
	//Services & configuration
	manager     service.Manager
	namespace   string
	name        string
	theme       ThemeProvider
	debounceTag int

	// Secret Data
	certPaginator  paginator.Model
	certViewPages  []string
	selectedSecret *secretsListItem
	secretsList    list.Model
	tlsSecrets     []tls.Secret

	// Ui elements
	selectedPane      Pane
	showRaw           bool
	loading           bool
	inspectedError    error
	errorModalMsg     string
	helpView          HelpViewModel
	spinner           spinner.Model
	inspectedViewport viewport.Model
	uiLayout          uiLayout
}

func NewModel(manager service.Manager, namespace, name string) (Model, error) {
	var items []list.Item
	secretsList := list.New(items, newSecretDelegate(), 50, 20)
	secretsList.Title = "Select a Secret Secret"
	secretsList.SetShowHelp(false)
	defaultPane := LeftPane
	return Model{
		certPaginator:     paginator.New(),
		inspectedViewport: viewport.New(50, 20), // Will be updated later,
		name:              name,
		namespace:         namespace,
		manager:           manager,
		secretsList:       secretsList,
		selectedPane:      defaultPane,
		spinner:           spinner.New(),
		theme:             Default,
		helpView:          NewHelpViewModel(defaultPane, Default),
	}, nil
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg { return loadSecretsMsg{} }
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
				m.certPaginator.PrevPage()
				m.inspectedViewport.SetContent(m.certViewPages[m.certPaginator.Page] + "\n\n" + m.certPaginator.View())
			case "right":
				m.certPaginator.NextPage()
				m.inspectedViewport.SetContent(m.certViewPages[m.certPaginator.Page] + "\n\n" + m.certPaginator.View())
			}
		}
		if m.secretsList.FilterState() != list.Filtering {
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
		m.handleCopyMsg(msg)
	case renderSecretsListMsg:
		//TODO here should be the filtering done
		m.secretsList.SetItems(tlsSecretsToListItems(m.tlsSecrets))
	case secretsLoadedMsg:
		m.tlsSecrets = msg.secrets
		m.loading = false
		cmds = append(cmds, func() tea.Msg { return renderSecretsListMsg{} })
	case switchCertViewMsg:
		m.showRaw = !m.showRaw
		cmds = append(cmds, func() tea.Msg { return inspectTLSSecretMsg{tag: m.debounceTag} })
	case switchPaneMsg:
		m.selectedPane = nextPane(m.selectedPane)
		m.helpView.SetPane(m.selectedPane)
	case loadSecretsMsg:
		cmds = append(cmds, loadSecretsCmd(m))
	case startLoadingMsg:
		m.loading = true
		cmds = append(cmds, m.spinner.Tick)
	case inspectTLSSecretMsg:
		if msg.tag == m.debounceTag {
			m.handleInspectTLSSecretMsg()
		}
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
			m.secretsList, listCmd = m.secretsList.Update(msg)
			cmds = append(cmds, listCmd)
		case RightPane:
			var vpCmd tea.Cmd
			m.inspectedViewport, vpCmd = m.inspectedViewport.Update(msg)
			cmds = append(cmds, vpCmd)
		}
	}

	if sel := m.secretsList.SelectedItem(); sel != nil {
		if item, ok := sel.(secretsListItem); ok {
			if m.selectedSecret == nil || !m.selectedSecret.Equals(item.Secret) {
				m.selectedSecret = &item
				m.debounceTag++
				cmds = append(cmds, tea.Tick(debounceDuration, func(t time.Time) tea.Msg { return inspectTLSSecretMsg{tag: m.debounceTag} }))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.errorModalMsg != "" {
		return m.renderErrorModal(m.errorModalMsg)
	}

	left := m.leftPane(m.uiLayout.LeftPaneWidth, m.uiLayout.UsableHeight)
	right := m.rightPane(m.uiLayout.RightPaneWidth, m.uiLayout.UsableHeight)

	mainContent := m.theme.DocStyle().Render(lipgloss.JoinHorizontal(lipgloss.Top, left, right))
	helpContent := m.helpView.View()

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, helpContent)
}

func loadSecretsCmd(m Model) tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return startLoadingMsg{} },
		func() tea.Msg {
			if m.manager == nil {
				return errorMsg{fmt.Errorf("secretsList service not initialized")}
			}

			if m.name != "" {
				tlsSecret, err := m.manager.LoadTLSSecret(m.namespace, m.name)
				if err != nil {
					return errorMsg{fmt.Errorf("failed to load secret %s/%s: %w", m.namespace, m.name, err)}
				}
				return secretsLoadedMsg{[]tls.Secret{tlsSecret}}
			}

			tlsSecrets, err := m.manager.ListTLSSecrets(m.namespace)
			if err != nil {
				return errorMsg{fmt.Errorf("failed to load secretsList in namespace %s: %w", m.namespace, err)}
			}

			return secretsLoadedMsg{tlsSecrets}
		},
	)
}

func (m Model) leftPane(width, height int) string {
	style := m.theme.Pane(m.selectedPane == LeftPane, width, height)
	if m.loading {
		return style.Render(m.spinner.View() + " Loading secretsList...")
	}
	return style.Render(m.secretsList.View())
}

func (m Model) rightPane(width, height int) string {
	style := m.theme.Pane(m.selectedPane == RightPane, width, height)
	if m.inspectedError != nil {
		return style.Render(fmt.Errorf("error inspecting secret: %w", m.inspectedError).Error())
	}
	if m.selectedSecret != nil && !m.loading {
		return style.Render(m.inspectedViewport.View())
	}
	return style.Render("Nothing yet selectedSecret, waiting.....")
}

func (m Model) renderErrorModal(msg string) string {
	errorModalRender := m.theme.ErrorModalWithWidth(m.uiLayout.TotalWidth).Render("Error:\n" + msg + "\n\nPress any key to dismiss.")
	return lipgloss.Place(m.uiLayout.TotalWidth, m.uiLayout.TotalHeight, lipgloss.Center, lipgloss.Center, errorModalRender)
}

func (m Model) inspectedTLSSecretContent() ([]string, error) {
	if m.showRaw {
		return []string{m.selectedSecret.PemCert(), m.selectedSecret.PemKey()}, nil
	}

	var views []string
	for _, cert := range m.selectedSecret.Certs() {
		views = append(views, formatCertificateInfo(cert, m.theme))
	}

	return views, nil
}

func (m *Model) updateLayout(width, height int) {
	m.uiLayout = calculateLayout(width, height, m.theme.DocStyle())
	m.secretsList.SetSize(m.uiLayout.LeftPaneWidth, m.uiLayout.UsableHeight)
	m.inspectedViewport.Width = m.uiLayout.RightPaneWidth
	m.inspectedViewport.Height = m.uiLayout.UsableHeight
	m.helpView.SetWidth(m.uiLayout.TotalWidth)
}

func (m *Model) handleInspectTLSSecretMsg() {
	if m.selectedSecret == nil {
		return
	}
	data, err := m.inspectedTLSSecretContent()
	m.certViewPages = data
	m.inspectedError = err
	m.certPaginator.SetTotalPages(len(data))
	m.certPaginator.Page = 0
	m.inspectedViewport.SetContent(m.certViewPages[m.certPaginator.Page] + "\n\n" + m.certPaginator.View())
}

func (m *Model) handleCopyMsg(msg copyMsg) {
	if m.selectedSecret == nil {
		return
	}

	var copyData string
	tlsCert, tlsKey := m.selectedSecret.PemCert(), m.selectedSecret.PemKey()

	if msg.key {
		copyData = tlsKey
	} else {
		copyData = tlsCert
	}

	if err := clipboard.WriteAll(copyData); err != nil {
		m.errorModalMsg = fmt.Sprintf("Error copying secret: %v", err)
	}
}

func tlsSecretsToListItems(ts []tls.Secret) []list.Item {
	items := make([]list.Item, len(ts))
	for i, s := range ts {
		items[i] = secretsListItem{s}
	}
	return items
}

func newSecretDelegate() secretsListDelegate {
	delegate := list.NewDefaultDelegate()

	// TODO: use the theme styles
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("#00BFFF"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(lipgloss.Color("#5DADE2"))

	return secretsListDelegate{delegate}
}

func nextPane(currentPane Pane) Pane {
	if currentPane == LeftPane {
		return RightPane
	}
	return LeftPane
}
