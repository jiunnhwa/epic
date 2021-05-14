{{ template "layout.tpl" . }}
{{ define "css" }}
        <link rel="stylesheet" href="/static/css/current.css">
{{ end}}

{{ define "content" }}
        <h2>{{ .PageTitle }}</h2>
        <p> This is msg: {{ .Msg }}</p>
{{ end }}

{{ define "js" }}
    <script src="/static/js/current.js"></script>
{{ end}}