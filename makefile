# SHELL := /bin/bash


# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

db-run:
	docker run --name smart-alert-pq \
		-e POSTGRES_USER=suser \
		-e POSTGRES_PASSWORD=spassword \
		-e POSTGRES_DB=smart-alert \
		-p 5432:5432 \
		-v "$(shell pwd)/.data/postgres:/var/lib/postgresql/data" \
		-d postgres:12-alpine

db-stop:
	docker stop smart-alert-pq
	docker rm smart-alert-pq

grafana-start:
	docker run -d \
  		-p 3000:3000 \
  		--name=smart-alert-grafana \
  		-e "GF_INSTALL_PLUGINS=grafana-clock-panel,grafana-simple-json-datasource" \
  		grafana/grafana:8.5.3

data-reset:
	PGPASSWORD=spassword psql -h localhost -U suser -d smart-alert -c "DELETE FROM datasets WHERE set = 'knc'; DELETE FROM alerts WHERE set = 'knc';"

data-gen:
	go run app/services/collector/main.go | go run app/tooling/logfmt/main.go -service=smart-alert-collector
	
monitor-run:
	go run app/services/alerter/main.go | go run app/tooling/logfmt/main.go -service=smart-alert-monitor