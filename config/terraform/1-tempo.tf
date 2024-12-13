# Install manually
# helm repo add grafana https://grafana.github.io/helm-charts
# helm repo update
# helm install tempo --namespace fitme-app-dev --create-namespace --version 1.6.2 --values terraform/values/tempo.yaml grafana/tempo
resource "helm_release" "tempo" {
  name = "tempo"

  repository       = "https://grafana.github.io/helm-charts"
  chart            = "tempo"
  namespace        = "fitme-app-dev"
  version          = "1.6.2"
  create_namespace = true

  values = [file("values/tempo.yaml")]
}
