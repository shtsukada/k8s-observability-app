{{- define "grpc-app.name" -}}
grpc-app
{{- end }}

{{- define "grpc-app.fullname" -}}
{{ include "grpc-app.name" . }}
{{- end }}