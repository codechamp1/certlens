package ui

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/codechamp1/certlens/internal/service"
)

type CertViewModel struct {
	certViews         []string
	certPages         paginator.Model
	inspectedViewport viewport.Model
	inspectRaw        bool
	inspectedError    error
	secretService     service.SecretsService
	theme             ThemeProvider
}

type switchCertViewMsg struct{}

func NewCertViewModel(svc service.SecretsService, theme ThemeProvider) CertViewModel {
	return CertViewModel{
		certPages:         paginator.New(),
		inspectedViewport: viewport.New(50, 20), // Will be updated later
		secretService:     svc,
		theme:             theme,
	}
}

func (c CertViewModel) Init() tea.Cmd {
	return nil
}

func (c CertViewModel) Update(msg tea.Msg) (CertViewModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			c.certPages.PrevPage()
			c.inspectedViewport.SetContent(c.certViews[c.certPages.Page] + "\n\n" + c.certPages.View())
		case "right":
			c.certPages.NextPage()
			c.inspectedViewport.SetContent(c.certViews[c.certPages.Page] + "\n\n" + c.certPages.View())
		}
	case switchCertViewMsg:
		c.inspectRaw = !c.inspectRaw
		cmds = append(cmds, func() tea.Msg { return inspectTLSSecretMsg{} })
	case inspectTLSSecretMsg:
		data, err := c.inspectedTLSSecretContent(msg.namespace, msg.name, c.inspectRaw)
		c.certViews = data
		c.inspectedError = err
		c.certPages.SetTotalPages(len(data))
		c.certPages.Page = 0
		c.inspectedViewport.SetContent(c.certViews[c.certPages.Page] + "\n\n" + c.certPages.View())
	}

	var vpCmd tea.Cmd
	c.inspectedViewport, vpCmd = c.inspectedViewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return c, tea.Batch(cmds...)
}

func (c CertViewModel) View() string {
	if c.inspectedError != nil {
		return fmt.Sprintf("error inspecting secret: %s", c.inspectedError)
	}

	return c.inspectedViewport.View()
}

type CertField struct {
	Label string
	Value string
}

func viewFieldsFromStruct(s interface{}) []CertField {
	var fields []CertField
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		label := fieldType.Tag.Get("label")
		if label == "" {
			continue
		}

		var strVal string
		switch fv := fieldValue.Interface().(type) {
		case []string:
			strVal = strings.Join(fv, ", ")
		case time.Duration:
			strVal = fv.String()
		case float64:
			strVal = fmt.Sprintf("%.2f", fv)
		default:
			strVal = fmt.Sprintf("%v", fv)
		}

		fields = append(fields, CertField{
			Label: label,
			Value: strVal,
		})
	}

	return fields
}

func renderField(keyStyle lipgloss.Style, valueStyle lipgloss.Style, key, value string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		keyStyle.Render(key+":"),
		valueStyle.Render(value),
	)
}

func formatCertificateInfo(ci service.CertificateInfo, t ThemeProvider) string {
	var sb strings.Builder
	val := reflect.ValueOf(ci)
	typ := reflect.TypeOf(ci)

	for i := 0; i < typ.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		label := fieldType.Tag.Get("label")
		if label == "" {
			continue
		}

		sb.WriteString(t.SectionHeader().Render(label))
		sb.WriteString("\n")
		fields := viewFieldsFromStruct(fieldVal.Interface())
		for _, f := range fields {
			sb.WriteString(renderField(t.Key(), t.Value(), f.Label, f.Value))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (c CertViewModel) inspectedTLSSecretContent(namespace, name string, raw bool) ([]string, error) {
	if raw {
		tlsCert, tlsKey, err := c.secretService.RawInspectTLSSecret(namespace, name)
		if err != nil {
			return nil, fmt.Errorf("failed to inspect secret %s/%s: %w", namespace, name, err)
		}
		return []string{tlsCert, tlsKey}, nil // o singură pagină
	}

	certs, err := c.secretService.InspectTLSSecret(namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect secret %s/%s: %w", namespace, name, err)
	}

	var views []string
	for _, cert := range certs {
		views = append(views, formatCertificateInfo(cert, c.theme))
	}
	return views, nil
}
