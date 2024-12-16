resource "helm_release" "otel_operator" {
  name       = "otel-operator"
  namespace  = "monitoring"
  chart      = "open-telemetry/opentelemetry-operator"
  repository = "otel/opentelemetry-collector-k8s"
  version    = "0.75.1"
  values = [
    file("values/otel-operator.yaml")
  ]
}
