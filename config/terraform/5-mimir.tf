# helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
# helm repo update
#
# helm install prometheus prometheus-community/prometheus \
# --namespace fitmeapp \
# --create-namespace --values terraform/values/prometheus.yaml

#helm install prometheus prometheus-community/kube-prometheus-stack --version "66.5.0" --namespace monitoring

resource "helm_release" "mimir" {
  name             = "mimir"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "mimir-distributed"
  namespace        = "monitoring"
  version          = "5.5.1"
  create_namespace = true
  values = [file("values/mimir.yaml")]
}
