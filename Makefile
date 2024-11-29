VERSION ?= latest
PREV_VERSION ?= 0.1.4
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

build-image:
	docker build --no-cache --platform=linux/amd64 -t fit-me:$(VERSION) -f Dockerfile .
	docker build --no-cache --platform=linux/amd64 -t fit-me:$(PREV_VERSION) -f Dockerfile .


tag-image:
	docker tag fit-me:$(VERSION) a11199/fit-me:latest
	docker tag fit-me:$(PREV_VERSION) a11199/fit-me:$(PREV_VERSION)

push-image:
	docker push a11199/fit-me:latest
	docker push a11199/fit-me:$(PREV_VERSION)
