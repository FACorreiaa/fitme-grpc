resource "helm_release" "loki" {
  name       = "loki"

  repository = "https://grafana.github.io/helm-charts"
  chart      = "loki-stack"
  namespace  = "monitoring"
  version    = "2.10.2"
  create_namespace = true

  values = [file("values/loki.yaml")]
}
