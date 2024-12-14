resource "helm_release" "postgres-exporter" {
  name             = "postgres-exporter"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "prometheus-postgres-exporter "
  namespace        = "monitoring"
  version          = "6.7.1"
  create_namespace = true
  values = [file("values/postgres-exporter.yaml")]
}
