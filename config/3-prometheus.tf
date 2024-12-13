# helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
# helm repo update
#
# helm install prometheus prometheus-community/prometheus \
# --namespace fitme-app-dev \
# --create-namespace --values terraform/values/prometheus.yaml
resource "helm_release" "prometheus" {
  name             = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "prometheus"
  namespace        = "fitme-app-dev"
  version          = "15.11.1"
  create_namespace = true

  values = [file("values/prometheus.yaml")] # Define your Prometheus configuration here
}
