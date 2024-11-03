# fitme-grpc

Rewrite of the FitME rest version to a  gRPC infrastructure with a strong emphasis on middleware, tracing, and logging integration.
The use of a shared proto directory and containerized service architecture gives plenty of flexibility in handling upstream services and dependencies.

The approach to middleware management—ensuring correct ordering for context propagation, handling panics, and managing logs and metrics, wrapping OpenTelemetry's middleware for future-proofing against breaking changes.

TODO
 - Finish all services
 - Message system to communicate between PT and its clients (add communications between Institution and PTs?)
 - Kafka message queue system and notifications between users and PTs
 - Fix Prometheus
 - Loki
 - Complete Grafana
 - PDF builder: https://github.com/johnfercher/maroto
 - CSV and Excel builder: check Domonda
 - Complete:
   - Traces
      Exporters: Stdout, Jaeger, Zipkin, Datadog and OpenTelemetry (OTLP) collector
      Importers: OpenTracingShim
     - Metrics
        Exporters: Prometheus, Datadog, and OpenTelemetry (OTLP) collector
        Importers: SwiftMetricsShim
       - Logs
        Exporters: OpenTelemetry (OTLP) collector
 - K8s deployment
 - Where to deploy ?
