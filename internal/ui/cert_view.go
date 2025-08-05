package ui

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/codechamp1/certlens/internal/domains/cert"
)

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

func formatCertificateInfo(ci cert.TLS, t ThemeProvider) string {
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
