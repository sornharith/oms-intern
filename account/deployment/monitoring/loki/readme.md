 helm install --values values.yaml loki --namespace=monitoring grafana/loki

send log to -> http://loki-gateway.monitoring.svc.cluster.local/loki/api/v1/push
