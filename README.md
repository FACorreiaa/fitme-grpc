# fitme-grpc

Rewrite of the FitME rest version to a  gRPC infrastructure with a strong emphasis on middleware, tracing, and logging integration.
The use of a shared proto directory and containerized service architecture gives plenty of flexibility in handling upstream services and dependencies.

The approach to middleware management—ensuring correct ordering for context propagation, handling panics, and managing logs and metrics, wrapping OpenTelemetry's middleware for future-proofing against breaking changes.

TODO
 - Fix Prometheus
 - Loki
 - Tempo
 - Complete Grafana (In production, configure ingres and point to prometheus-grafana)
 - All above with Kubernetes
 - Finish all services
 - Message system to communicate between PT and its clients (add communications between Institution and PTs?)
 - Kafka (?) message queue system and notifications between users and PTs
 - PDF builder: https://github.com/johnfercher/maroto
 - CSV and Excel builder: check Domonda
- K8s deployment
- Where to deploy ?


## Traces
- Exporters: Stdout, Jaeger, Zipkin, Datadog and OpenTelemetry (OTLP) collector
- Importers: OpenTracingShim
## Metrics
- Exporters: Prometheus, Datadog, and OpenTelemetry (OTLP) collector
- Importers: SwiftMetricsShim
## Logs
- Exporters: OpenTelemetry (OTLP) collector

For dev:
Tempo => kubectl port-forward svc/tempo 4317 -n monitoring
Grafana => kubectl port-forward svc/grafana 3000:80 -n monitoring
Prometheus => kubectl port-forward prometheus-prometheus-kube-prometheus-prometheus-0 9090 -n monitoring
