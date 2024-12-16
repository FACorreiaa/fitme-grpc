resource "helm_release" "otel_collector" {
  name       = "otel-collector"
  namespace  = "monitoring"
  chart      = "open-telemetry/opentelemetry-collector"
  version    = "0.110.7"
  values = [
    file("values/otel-collector-helm.yaml")
  ]
}
