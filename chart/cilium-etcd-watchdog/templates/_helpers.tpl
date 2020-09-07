{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "cilium-etcd-watchdog.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* vim: set filetype=mustache: */}}
{{/*
Expand the namespace of the chart.
*/}}
{{- define "cilium-etcd-watchdog.namespace" -}}
{{- default "kube-system" .Values.namespaceOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cilium-etcd-watchdog.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cilium-etcd-watchdog.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "cilium-etcd-watchdog.labels" -}}
helm.sh/chart: {{ include "cilium-etcd-watchdog.chart" . }}
{{ include "cilium-etcd-watchdog.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "cilium-etcd-watchdog.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cilium-etcd-watchdog.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
name: cilium-etcd-operator  # https://docs.cilium.io/en/v1.8/concepts/terminology/#well-known-identities
io.cilium/app: etcd-operator  # https://docs.cilium.io/en/v1.8/concepts/terminology/#well-known-identities
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "cilium-etcd-watchdog.serviceAccountName" -}}
cilium-etcd-operator  # https://docs.cilium.io/en/v1.8/concepts/terminology/#well-known-identities
{{- end -}}

{{/*
Expand the container image.
*/}}
{{- define "cilium-etcd-watchdog.image" -}}
{{- $tag := default .Chart.AppVersion .Values.image.tag -}}
{{- printf "%s:%s" .Values.image.repository $tag -}}
{{- end -}}
