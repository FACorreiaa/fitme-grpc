# helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
# helm repo update
#
# helm install prometheus prometheus-community/prometheus \
# --namespace fitmeapp \
# --create-namespace --values terraform/values/prometheus.yaml

#helm install prometheus prometheus-community/kube-prometheus-stack --version "66.5.0" --namespace fitmeapp

resource "helm_release" "prometheus" {
  name             = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  namespace        = "fitmeapp"
  version          = "66.5.0"
  create_namespace = true

  values = [
    yamlencode({
      additionalServiceMonitors = [{
        name      = "fitme-service-monitor"
        namespace = "fitmeapp"
        selector  = {
          matchLabels = {
            app = "fitme"
          }
        }

        endpoints = [{
          port     = "http"
          interval = "30s"
          path     = "/metrics"
        }]
      }]
    })
  ]
}
