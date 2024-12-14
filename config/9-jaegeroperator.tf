resource "helm_release" "jaeger-operator" {
  name             = "jaegertracing"
  repository       = "https://jaegertracing.github.io/helm-charts"
  chart            = "jaeger-operator "
  namespace        = "monitoring"
  version          = "2.57.0"
  create_namespace = true
  values = [file("values/jaeger-operator.yaml")]
}
