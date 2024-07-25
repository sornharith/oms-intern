helm install loki-distributed grafana/loki-distributed -f values.yaml -n monitoring

send log to -> http://loki-gateway.monitoring.svc.cluster.local/loki/api/v1/push
