# fitme-grpc

Rewrite of the FitME rest version to a  gRPC infrastructure with a strong emphasis on middleware, tracing, and logging integration.
The use of a shared proto directory and containerized service architecture gives plenty of flexibility in handling upstream services and dependencies.

The approach to middleware management—ensuring correct ordering for context propagation, handling panics, and managing logs and metrics, wrapping OpenTelemetry's middleware for future-proofing against breaking changes.

TODO
 - Fix Prometheus [x]
 - Loki [x]
 - Tempo [x]
 - Complete Grafana (In production, configure ingres and point to prometheus-grafana)
 - All above with Kubernetes [x]
 - Finish postgres and redis exporter [x]
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

Docker compose just for local testing, telemetry services not relevant in there but all working within Kubernetes.
For dev:
Tempo => kubectl port-forward svc/tempo 4317 -n monitoring
Grafana => kubectl port-forward svc/grafana 3000:80 -n monitoring
Prometheus => kubectl port-forward prometheus-prometheus-kube-prometheus-prometheus-0 9090 -n monitoring

# Leaderboard feature
=> I am thinking about opening the user platform so every user can see other users progress and plans made etc and introducing a leaderboard system with points for achievements completed. :aPES_Think:
This only for regular users. A "PT" would still only manage their own clients and not see other "PT" clients, obviously.

=> maybe you want progress to be anonymized too. so people can only see the progress in terms of % towards some goal

=> i'm not too sure how embarrassed people are about their fitness goals. but it's one of the main things that keeps them out of gyms from what i can tell reading on the internet

=> :pog: progress leaderboard for friends who train together. can of worms to get into, you want the leaderboards to
have customizable goals (absolute/relative/custom) be public/private  look cute
