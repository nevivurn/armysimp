package report

import (
	"bytes"
	"text/template"
	"time"
)

const (
	templateString = `armysimp report
{{ .Timestamp.Format "Mon, 02 Jan 2006 15:04:05 MST" }}
{{ range .Generations }}
{{ .Title }}:
{{- range .Channels }}
- {{ .Title }} {{ .Subscribers }}
{{- range .Videos }}
+ {{ .TitleTranslated }}
{{- end }}
{{- end }}
{{- end }}
`

	templateStringFavorites = `armysimp report
{{ .Timestamp.Format "Mon, 02 Jan 2006 15:04:05 MST" }}
{{ range .Generations }}
{{ .Title }}:
{{- range .Channels }}
{{- if .Favorite }}
- {{ .Title }} {{ .Subscribers }}
{{- range .Videos }}
+ {{ .TitleTranslated }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
`
)

var (
	reportTemplate          = template.Must(template.New("report").Parse(templateString))
	reportTemplateFavorites = template.Must(template.New("report").Parse(templateStringFavorites))
)

// Data is the input go generate the report.
type Data struct {
	Timestamp   time.Time
	Generations []Generation
}

// Generation is a collection of YT channels.
type Generation struct {
	Title    string
	Channels []Channel
}

// Channel is a single channel.
type Channel struct {
	Title       string
	Favorite    bool
	Subscribers uint64

	Videos []Video
}

// Video is a single video in a channel.
type Video struct {
	Title           string
	TitleTranslated string
}

// Generate a report.
func Generate(data Data) (string, error) {
	buf := bytes.Buffer{}
	err := reportTemplate.Execute(&buf, data)
	return buf.String(), err
}

// GenerateFavorites generates a favorites report.
func GenerateFavorites(data Data) (string, error) {
	buf := bytes.Buffer{}
	err := reportTemplateFavorites.Execute(&buf, data)
	return buf.String(), err
}
