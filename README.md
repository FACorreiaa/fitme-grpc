# FitME-gRPC

Rewrite of the FitME REST version to a gRPC infrastructure with a strong emphasis on middleware, tracing, and logging integration.
By sharing proto definitions and containerizing services, this project achieves flexible inter-service communication while ensuring proper context propagation, error handling, and telemetry across the stack.

## Overview
FitME-gRPC is a rewrite of the original FitME REST API into a gRPC-based service architecture. The project is built with a strong emphasis on robust middleware management, comprehensive tracing and logging, and seamless integration with modern observability tools (Prometheus, Grafana, Loki, Tempo). It aims to improve scalability and maintainability while providing enhanced telemetry and debugging capabilities.

## Architecture & Technology Stack

- Backend:
1. Language: Go (Golang)
2. Communication: gRPC with shared proto definitions
3. Middleware: Custom middleware layers wrapping OpenTelemetry for tracing, logging, and metrics
4. Database: PostgreSQL
5. Caching: Redis

- Telemetry:
1. Tracing: Stdout, Jaeger, Zipkin, Datadog, and OTLP collector
2. Metrics: Prometheus, Datadog, and OTLP collector
3. Logs: OTLP collector (integrated with Loki)
4. Containerization: Docker (with Kubernetes deployment for production)

- Other Tools:
1. PDF builder: Maroto
2. CSV/Excel builder: (see Domonda for inspiration)

## Core Features
- gRPC Infrastructure:
Fully containerized and using gRPC for high-performance communication between services.

- Middleware Management:
Ensures proper ordering for context propagation, panic handling, logging, and metrics collection. Uses wrappers around OpenTelemetry middleware to guard against breaking changes.

- Observability:
Integrated with multiple telemetry backends including Prometheus, Grafana, Loki, Tempo, and Jaeger.

- Data Exporters:
Support for PostgreSQL, Redis exporters, and integration with message systems for notifications (e.g., Kafka).

## Telemetry
### Traces
- Exporters:
Stdout
Jaeger
Zipkin
Datadog
OpenTelemetry (OTLP) Collector
- Importers:
OpenTracingShim

### Metrics
- Exporters:
Prometheus
Datadog
OpenTelemetry (OTLP) Collector
- Importers:
SwiftMetricsShim
Logs
Exporters:
OpenTelemetry (OTLP) Collector
(Integrated with Loki for log aggregation)

## Service Integration
- Proto Sharing:
A shared proto directory is used for service communication across the architecture.
- Containerized Services:
Docker Compose is used for local testing, while Kubernetes is used in production for scalability and high availability.
- Middleware Pipeline:
Middleware ensures correct ordering of interceptors for tracing, logging, and metrics (e.g., Prometheus and OTEL interceptors).

## Deployment & Kubernetes
Local Testing:

- Tempo:

```kubectl port-forward svc/tempo 4317 -n monitoring```
- Grafana:
```kubectl port-forward svc/grafana 3000:80 -n monitoring```
Prometheus:
```kubectl port-forward prometheus-prometheus-kube-prometheus-prometheus-0 9090 -n monitoring```
Production Deployment:
- All telemetry services (Prometheus, Grafana, Loki, Tempo, etc.) are deployed on Kubernetes.
- Ingress configuration should be updated in production to point to the proper Grafana/Prometheus endpoints.

### Additional Features
- Leaderboard Feature
Concept:
Open the user platform to allow regular users to see progress, plans, and achievements.
- Privacy:
Consider anonymizing data and providing customizable, public or private leaderboards.
- Customization:
Allow users to define goals (absolute, relative, custom) and compare progress with friends or community members.

## FitSynch
FitSynch expands the FitME concept by integrating personal training management with AI-powered meal planning and shopping assistance.
Key Components:

- User Management for Trainers:
Manage clients, assign workout/meal plans, and communicate via messaging.
- Workout & Meal Plans:
Create personalized workout plans with video tutorials and customizable meal plans with macro breakdowns.
- Ingredients, Recipes, and Shopping List Generator:
Build a database of nutritional data, share recipes, and automatically generate shopping lists.
- AI-Powered Assistance:
1. Meal Plan AI: Suggest meals based on preferences and dietary restrictions.
2. Shopping AI: Provide calorie breakdowns, healthier alternatives, and cost-saving suggestions.
3. Fitness Insights: Analyze user progress and recommend workout adjustments.
- Trainer Dashboard:
A centralized dashboard for trainers to manage clients, view progress, and handle payments.

### Meal Plan Validation
- Purpose:
Ensure that a meal plan aligns with the user's objective (e.g., maintenance mode should not exceed a specific calorie goal).
- Guidance:
Warn users if their meal plan exceeds their objective's calorie goal, allowing for adjustments or confirmations.
- Flexibility:
Optionally allow overrides while logging such events for further review.

- TODO / Future Work

[x] Fix Prometheus integration
[x]Configure Loki integration
[x]Configure Tempo integration
[x]Complete Grafana configuration for production (configure ingress and point to Prometheus)
[x]Kubernetes deployment for all services
[x]Finalize PostgreSQL and Redis exporters
[]Finalize all remaining services
[]Implement a messaging system for communication between personal trainers (PTs) and clients
[]Integrate Kafka (or similar) for message queue and notifications
[]Add PDF builder (using Maroto)
[]Add CSV and Excel builders (see Domonda for reference)
[]Further enhance security, data privacy, and RBAC
[]Finalize the leaderboard feature and FitSynch enhancements

### Access Telemetry UIs:

Grafana: http://localhost:3000
Jaeger: http://localhost:16686
Prometheus: http://localhost:9090
Loki: Use the mapped port (e.g., http://localhost:3100 for API requests)
Kubernetes (Production):
Follow the Kubernetes deployment manifests and use port-forwarding as outlined in the Deployment & Kubernetes section for accessing services locally.

Contributing
Contributions are welcome! Please fork the repository, create your feature branch, and submit a pull request. Ensure that your changes are covered by appropriate tests and documentation updates.
