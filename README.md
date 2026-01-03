This is a comprehensive overhaul of your **FitME-gRPC** documentation. I’ve structured it to be professional, scannable, and ready for both developers and potential stakeholders (like HIMS & HERS partners).

---

# 🚀 FitME-gRPC

### The High-Performance, AI-Driven Wellness Ecosystem

**FitME-gRPC** is an enterprise-grade rewrite of the original FitME REST API. Built with **Go** and **gRPC**, it transitions from a traditional monolithic structure to a highly observable, microservice-based architecture. It integrates fitness tracking, AI nutritional coaching, mental health journaling, and biometric data synchronization into a single, unified platform.

---

## 🏗 Service Layer Refactor (ASAP Priority)

To ensure a clean separation of concerns, the service layer is currently being refactored to follow the **Hexagonal Architecture** pattern:

* **gRPC Interface:** Handles incoming proto-generated requests.
* **Domain Logic:** Core business rules (converts proto messages to domain entities).
* **Orchestration:** Coordinates between repositories (DB) and external microservices.
* **Database Schema:** Strictly limited to persistence logic and data integrity.

---

## 🛠 Technology Stack

| Layer | Technology |
| --- | --- |
| **Backend** | Go (Golang), gRPC, Protobuf |
| **Database** | PostgreSQL (TimescaleDB for biometrics) |
| **Caching** | Redis (Real-time metrics & session management) |
| **Messaging** | Kafka (Event-driven architecture) |
| **Observability** | Prometheus, Grafana, Loki, Tempo, Jaeger |
| **Deployment** | Docker, Kubernetes (K8s) |
| **Integrations** | Apple Health, Oura, Whoop, FHIR Standards |

---

## 🌐 System Architecture

The architecture is designed for high availability and deep observability. Every service is wrapped in custom middleware for tracing and telemetry.

### Core Microservices

1. **UserService:** Identity, RBAC, and social graph.
2. **BiometricsService:** Real-time ingestion of wearable data (HRV, Sleep, SpO2).
3. **AIOrchestrationService:** Generates personalized meal plans and injury risk assessments.
4. **MentalHealthService:** Encrypted journaling and AI-driven mood analysis.
5. **DietService:** Macro tracking and automated shopping list generation.
6. **TrainerService:** Managing PT-client relationships and billing.

---

## 📈 Observability & Telemetry

We treat observability as a first-class citizen. The stack ensures that every request can be traced from the mobile app through the gateway to the database.

* **Traces:** Distributed tracing via **Tempo** and **Jaeger**.
* **Metrics:** Real-time performance monitoring via **Prometheus**.
* **Logs:** Aggregated log streams in **Loki**.
* **Dashboards:** Unified visualization in **Grafana**.

---

## 🗺 Roadmap & Milestones

### Phase 1: Foundation (Current)

* [x] Migrate REST endpoints to gRPC.
* [x] Implement OpenTelemetry middleware.
* [x] Establish K8s deployment manifests.
* [ ] Finalize PostgreSQL/Redis exporters.

### Phase 2: Intelligence & Social

* [ ] **FitSynch:** AI-powered meal planning and shopping list generator.
* [ ] **Leaderboards:** Public/Private competitive gamification.
* [ ] **Messaging:** Bi-directional streaming for PT-Client chat.

### Phase 3: Healthcare Integration

* [ ] **FHIR Compliance:** Implement SMART on FHIR for medical data portability.
* [ ] **HIMS & HERS Integration:** Telehealth partnership for clinical insights.
* [ ] **Predictive Analytics:** AI-based burnout and depression risk detection.

---

## 🧬 Data Flow: Biometric Ingestion

The following sequence illustrates how wearable data (like a heart rate spike) moves from a device to a clinical insight.

---

## 🔐 Security & Compliance

Given the sensitive nature of health data, FitME-gRPC follows strict security protocols:

* **End-to-End Encryption:** For all mental health journaling and private messages.
* **HIPAA Preparedness:** Data residency and anonymization policies for research.
* **JWT/OAuth2:** Secure service-to-service communication with scoped permissions.

---

## 🤝 Technical Partnership: HIMS & HERS

FitME is uniquely positioned to bridge the gap between daily fitness and clinical telehealth. By implementing **FHIR (Fast Healthcare Interoperability Resources)**, we allow healthcare providers on the HIMS platform to view real-time patient recovery and biometric trends.

> **Key Integration Benefit:** "Moving from reactive medicine to proactive wellness by providing doctors with 24/7 biometric context."

---

## 🚀 Getting Started

### Local Development (Docker Compose)

1. **Clone the repo:** `git clone https://github.com/your-repo/fitme-grpc`
2. **Start the stack:** `docker-compose up -d`
3. **Access UIs:**
* Grafana: `http://localhost:3000`
* Jaeger: `http://localhost:16686`
* Prometheus: `http://localhost:9090`



### Kubernetes (Production)

```bash
# Deploy monitoring stack
kubectl apply -k ./deploy/k8s/monitoring

# Deploy core services
kubectl apply -k ./deploy/k8s/services

```

---

## ✉️ Contributing

We welcome contributions! Please see our `CONTRIBUTING.md` for coding standards, specifically regarding proto definitions and domain-driven design patterns.

**Would you like me to refine the `.proto` file definitions for the `AIOrchestrationService` to match this new architecture?**