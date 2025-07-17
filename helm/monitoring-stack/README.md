# Monitoring Stack Helm Chart

このディレクトリに「kube-prometheus-stack」のカスタム"values.yaml"を格納。
https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack

Prometheus, Grafana, Alertmanager, Node Exporter, Kube State Metrics をまとめて導入

※事前準備

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
kubectl create namespace monitoring