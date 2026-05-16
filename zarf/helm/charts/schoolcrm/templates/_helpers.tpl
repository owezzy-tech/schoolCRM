{{/*
Common labels
*/}}
{{- define "schoolcrm.labels" -}}
app: schoolcrm
app.kubernetes.io/name: schoolcrm
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "schoolcrm.selectorLabels" -}}
app: schoolcrm
app.kubernetes.io/name: schoolcrm
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Database environment variables for SchoolCRM service
*/}}
{{- define "schoolcrm.dbEnvVars" -}}
- name: SCHOOLCRM_DB_USER
  valueFrom:
    configMapKeyRef:
      name: schoolcrm-config
      key: db_user
- name: SCHOOLCRM_DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: schoolcrm-secret
      key: db_password
- name: SCHOOLCRM_DB_HOST_PORT
  valueFrom:
    configMapKeyRef:
      name: schoolcrm-config
      key: db_hostport
- name: SCHOOLCRM_DB_DISABLE_TLS
  valueFrom:
    configMapKeyRef:
      name: schoolcrm-config
      key: db_disabletls
{{- end -}}

{{/*
Kubernetes metadata environment variables
*/}}
{{- define "schoolcrm.k8sEnvVars" -}}
- name: KUBERNETES_NAMESPACE
  valueFrom:
    fieldRef:
      fieldPath: metadata.namespace
- name: KUBERNETES_NAME
  valueFrom:
    fieldRef:
      fieldPath: metadata.name
- name: KUBERNETES_POD_IP
  valueFrom:
    fieldRef:
      fieldPath: status.podIP
- name: KUBERNETES_NODE_NAME
  valueFrom:
    fieldRef:
      fieldPath: spec.nodeName
{{- end -}}
