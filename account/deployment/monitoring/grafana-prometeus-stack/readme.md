
for install  >  ``helm -n monitoring install promet-grafana prometheus-community/kube-prometheus-stack -f values.yaml``
<br>
for update > ``helm upgrade my-grafana grafana/grafana -f values.yaml -n monitoring``