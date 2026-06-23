{{/*
Раскрывает имя чарта.
*/}}
{{- define "bookings.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Создает полное имя приложения по умолчанию.
Обрезаем до 63 символов, потому что некоторые поля имен Kubernetes ограничены этим размером по спецификации DNS.
Если имя релиза содержит имя чарта, оно будет использовано как полное имя.
*/}}
{{- define "bookings.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Создает имя и версию чарта для метки чарта.
*/}}
{{- define "bookings.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Общие метки
*/}}
{{- define "bookings.labels" -}}
helm.sh/chart: {{ include "bookings.chart" . }}
{{ include "bookings.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Метки селектора
*/}}
{{- define "bookings.selectorLabels" -}}
app.kubernetes.io/name: {{ include "bookings.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Создает имя service account для использования
*/}}
{{- define "bookings.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "bookings.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
