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
	go build -o pprofbin .
	go tool pprof -http=":6969" pprofbin ./cpu.pprof
