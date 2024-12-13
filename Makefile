VERSION ?= latest
PREV_VERSION ?= 0.1.5
image_name = fit-me


run-down:
	docker compose down

run-up:
	docker compose up -d

restart-db:
	docker compose down && rm -rf ./.data && docker compose up -d

log-p:
	docker logs --details --follow --timestamps --tail=1000 inkme-dev-postgres

log-r:
	docker logs --details --follow --timestamps --tail=1000 inkme-dev-redis

run-prom:
	prometheus --config.file=config/prometheus.yml

go-lint: ## Runs linter for .go files
	@golangci-lint run --config ./config/go.yml
	@echo "Go lint passed successfully"

go-pprof:
	go tool pprof http://localhost:6060/debug/pprof/profile

update:
	go get -u

down-dev:
	docker compose down
	rm -rf .data

run-test:
	go test ./...

test-lint:
	testifylint --fix ./...

profile:
	go tool pprof \
      -raw -output=cpu.txt \
      'http://localhost:8080/debug/pprof/profile?seconds=60'

profile-graph:
	stackcollapse-go.pl cpu.txt | flamegraph.pl > cpu.svg

# --platform=linux/amd64
build-image:
	docker build --no-cache -t fit-me:$(VERSION) -f Dockerfile .
	docker build --no-cache -t fit-me:$(PREV_VERSION) -f Dockerfile .


tag-image:
	docker tag fit-me:$(VERSION) a11199/fit-me:latest
	docker tag fit-me:$(PREV_VERSION) a11199/fit-me:$(PREV_VERSION)

push-image:
	docker push a11199/fit-me:latest
	docker push a11199/fit-me:$(PREV_VERSION)

# --platform linux/amd64
#run-locally:
#	docker run -it --rm -p 8000:8000 -p 8001:8001 a11199/fit-me:0.1.4

run-locally:
	docker run -it --rm \
	  -p 8000:8000 -p 8001:8001 \
      -e POSTGRES_USER=postgres \
      -e POSTGRES_PASSWORD=postgres \
      -e POSTGRES_DB=fit-me-dev \
      -e POSTGRES_PORT=5440 \
      -e POSTGRES_HOST=postgres \
      -e POSTGRES_SCHEMA="" \
      -e REDIS_HOST="redis" \
      -e REDIS_PASSWORD="" \
      -e REDIS_PORT="6388" \
      -e JWT_SECRET_KEY="YOURMOMISHOT" \
      -e SCHEMA="" \
      -e CERT_FILE=./.data/server.crt \
      -e KEY_FILE=./.data/server.key \
      -e LOG_LEVEL=info \
      -e MODE=production \
      -e ADMIN_PASSWORD=elitahadmin!2024 \
      -e OTEL_SERVICE_NAME=FitMeDev \
      -e OTEL_RESOURCE_ATTRIBUTES="deployment.environment=production,service.namespace=Workouts,service.version=0.1,service.instance.id=localhost:8000" \
      -e OTEL_EXPORTER_OTLP_ENDPOINT=http://jaeger:4317/v1/traces \
      -e OTEL_EXPORTER_OTLP_PROTOCOL=grpc \
      -e OTEL_EXPORTER_INSECURE=true \
      -e OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=otel-collector:4317 \
      -e KUBERNETES_SERVICE_HOST="" \
      a11199/fit-me:debug


run-debug:
	docker run -it --rm --entrypoint sh \
      -e POSTGRES_USER=postgres \
      -e POSTGRES_PASSWORD=postgres \
      -e POSTGRES_DB=fit-me-dev \
      -e POSTGRES_PORT=5440 \
      -e POSTGRES_HOST=postgres \
      -e POSTGRES_SCHEMA="" \
      -e REDIS_HOST="redis" \
      -e REDIS_PASSWORD="" \
      -e REDIS_PORT="6388" \
      -e JWT_SECRET_KEY="YOURMOMISHOT" \
      -e SCHEMA="" \
      -e CERT_FILE=./.data/server.crt \
      -e KEY_FILE=./.data/server.key \
      -e LOG_LEVEL=info \
      -e MODE=production \
      -e ADMIN_PASSWORD=elitahadmin!2024 \
      -e OTEL_SERVICE_NAME=FitMeDev \
      -e OTEL_RESOURCE_ATTRIBUTES="deployment.environment=production,service.namespace=Workouts,service.version=0.1,service.instance.id=localhost:8000" \
      -e OTEL_EXPORTER_OTLP_ENDPOINT=http://jaeger:4317/v1/traces \
      -e OTEL_EXPORTER_OTLP_PROTOCOL=grpc \
      -e OTEL_EXPORTER_INSECURE=true \
      -e OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=otel-collector:4317 \
      a11199/fit-me:debug


run-debug-arm:
	docker run -it --rm --entrypoint sh a11199/fit-me:0.1.4

namespace:
	kubectl config set-context --current --namespace=fitmeapp

watch-tempo:
	kubectl port-forward svc/tempo 4317 -n monitoring

watch-grafana:
	kubectl port-forward svc/grafana 3000:80 -n monitoring

watch-prometheus:
	kubectl port-forward prometheus-prometheus-kube-prometheus-prometheus-0 9090 -n monitoring
