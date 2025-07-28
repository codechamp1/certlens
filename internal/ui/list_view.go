package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/codechamp1/certlens/internal/domains"
)

type ListViewModel struct {
	secrets list.Model
}

type secretsLoadedMsg struct{}

type updateSecretsMsg struct {
	secretsIDs []domains.K8SResourceID
}

type selectedSecretMsg struct {
	name      string
	namespace string
}

type secretItem struct {
	name      string
	namespace string
}

func (s secretItem) Title() string       { return s.name }
func (s secretItem) Description() string { return "Namespace: " + s.namespace }
func (s secretItem) FilterValue() string { return s.name }

func NewListViewModel() ListViewModel {
	var items []list.Item
	secretsList := list.New(items, newSecretDelegate(), 50, 20)
	secretsList.Title = "Select a TLS Secret"
	return ListViewModel{
		secrets: secretsList,
	}
}

func (lm ListViewModel) Init() tea.Cmd {
	return nil
}

func (lm ListViewModel) Update(msg tea.Msg) (ListViewModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case updateSecretsMsg:
		items := make([]list.Item, len(msg.secretsIDs))
		for i, s := range msg.secretsIDs {
			items[i] = secretItem{s.Name, s.Namespace}
		}
		lm.secrets.SetItems(items)
		cmds = append(cmds, func() tea.Msg { return secretsLoadedMsg{} })
	}

	var listCmd tea.Cmd
	lm.secrets, listCmd = lm.secrets.Update(msg)
	cmds = append(cmds, listCmd)

	if sel := lm.secrets.SelectedItem(); sel != nil {
		if item, ok := sel.(secretItem); ok {
			cmds = append(cmds, func() tea.Msg { return selectedSecretMsg{item.name, item.namespace} })
		}
	}

	return lm, tea.Batch(cmds...)
}

func (lm ListViewModel) View() string {
	return lm.secrets.View()
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
