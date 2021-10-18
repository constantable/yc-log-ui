{{/*
Expand the name of the chart.
*/}}
{{- define "yc-log-ui.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "yc-log-ui.fullname" -}}
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
Create chart name and version as used by the chart label.
*/}}
{{- define "yc-log-ui.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "yc-log-ui.labels" -}}
helm.sh/chart: {{ include "yc-log-ui.chart" . }}
{{ include "yc-log-ui.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "yc-log-ui.selectorLabels" -}}
app.kubernetes.io/name: {{ include "yc-log-ui.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "yc-log-ui.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "yc-log-ui.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}


{{- define ".env.variables" }}
          env:
            - name: SECRET
              value: {{ .Values.backend.parameters.secret | quote }}
            - name: DBFILE
              value: {{ .Values.backend.parameters.dbfile | quote }}
            - name: BACKEND_PORT
              value: {{ .Values.backend.parameters.backend_port | quote }}
            - name: SERVICE_ACCOUNT_ID
              value: {{ .Values.backend.parameters.service_account_id | quote }}
            - name: LOG_GROUP_ID
              value: {{ .Values.backend.parameters.log_group_id | quote }}
            - name: PUBLIC_KEY
              value: {{ .Values.backend.parameters.public_key | quote }}
            - name: PRIVATE_KEY
              value: {{ .Values.backend.parameters.private_key | quote }}
            - name: KEY_ID
              value: {{ .Values.backend.parameters.key_id | quote }}
{{- end }}