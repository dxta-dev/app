package markdown

import (
	"html/template"
)

var (
	AggregatedValuesTmpl = template.Must(template.New("aggregated_values").Parse(`
# {{ .Title }}

<description>
{{ .Description }}
</description>

**Overall:** {{ .Overall.Value }}

| Week       | Value |
|------------|-------|
{{- range .Weekly }}
| {{ .Week }} | {{ .Data.Value }} |
{{- end }}
`))
)
