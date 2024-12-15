resource "helm_release" "otel_collector" {
  name             = "otel-collector"
  namespace        = "monitoring"
  create_namespace = true
  repository       = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart            = "opentelemetry-collector"
  version          = "0.80.0"

  set {
    name  = "mode"
    value = "deployment" # Use "daemonset" if you want one collector per node
  }

  set {
    name  = "image.repository" # Key to set
    value = "otel/opentelemetry-collector-k8s" # Repository value
  }

  set {
    name = "command.name"
    value = "otelcol-k8s"
  }


}
