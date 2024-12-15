resource "helm_release" "jaeger-tracing" {
  name             = "jaegertracing"
  repository       = "https://jaegertracing.github.io/helm-charts"
  chart            = "jaeger"
  namespace        = "monitoring"
  version          = "3.3.3"
  create_namespace = true
  values = [file("values/jaeger.yaml")]
}
