# fitme-grpc

Rewrite of the FitME rest version to a  gRPC infrastructure with a strong emphasis on middleware, tracing, and logging integration.
The use of a shared proto directory and containerized service architecture gives plenty of flexibility in handling upstream services and dependencies.

The approach to middleware management—ensuring correct ordering for context propagation, handling panics, and managing logs and metrics, wrapping OpenTelemetry's middleware for future-proofing against breaking changes.
