package markdown

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/dxta-dev/app/internal/data"
)

var aggregatedValuesTmplText = `
# {{ .Title }}


## Description
{{ .Description }}

**Overall:** {{ .Values.Overall.Value }}

{{- $m := toMap .Values.Weekly -}}
| Week            | Value |
|-----------------|-------|
{{- range $week, $val := $m }}
| {{ $week }} | {{ $val }} |
{{- end }}
`

func aggregatedValuesToMap(wd []data.WeeklyData[data.Value]) map[string]int {
	m := make(map[string]int, len(wd))
	for _, w := range wd {
		if w.Data.Value != nil {
			m[w.Week] = *w.Data.Value
		}
	}
	return m
}

type aggregattedValuesPayload struct {
	Title       string
	Description string
	Values      *data.AggregatedValues
}

func GetAggregatedValuesMarkdown(
	ctx context.Context,
	title string,
	description string,
	values *data.AggregatedValues,
) (string, error) {
	tmpl, err := template.
		New("agg").
		Funcs(template.FuncMap{"toMap": aggregatedValuesToMap}).
		Parse(aggregatedValuesTmplText)
	if err != nil {
		return "", fmt.Errorf("parsing markdown template: %w", err)
	}

	var buf bytes.Buffer
	payload := aggregattedValuesPayload{
		Title:       title,
		Description: description,
		Values:      values,
	}
	if err := tmpl.Execute(&buf, payload); err != nil {
		return "", fmt.Errorf("executing markdown template: %w", err)
	}

	return buf.String(), nil
}
