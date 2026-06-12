# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# Deploy First Mentality

# ==============================================================================
# Go Installation
#
#	You need to have Go version 1.25 to run this code.
#
#	https://go.dev/dl/
#
#	If you are not allowed to update your Go frontend, you can install
#	and use a 1.25 frontend.
#
#	$ go install golang.org/dl/go1.25@latest
#	$ go1.25 download
#
#	This means you need to use `go1.25` instead of `go` for any command
#	using the Go frontend tooling from the makefile.

# ==============================================================================
# Brew Installation
#
#	Having brew installed will simplify the process of installing all the tooling.
#
#	Run this command to install brew on your machine. This works for Linux, Mac and Windows.
#	The script explains what it will do and then pauses before it does it.
#	$ /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
#
#	WINDOWS MACHINES
#	These are extra things you will most likely need to do after installing brew
#
# 	Run these three commands in your terminal to add Homebrew to your PATH:
# 	Replace <name> with your username.
#	$ echo '# Set PATH, MANPATH, etc., for Homebrew.' >> /home/<name>/.profile
#	$ echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> /home/<name>/.profile
#	$ eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
#
# 	Install Homebrew's dependencies:
#	$ sudo apt-get install build-essential
#
# 	Install GCC:
#	$ brew install gcc

# ==============================================================================
# Install Tooling and Dependencies
#
#	This project uses Docker and it is expected to be installed. Please provide
#	Docker at least 4 CPUs. To use Podman instead please alias Docker CLI to
#	Podman CLI or symlink the Docker socket to the Podman socket. More
#	information on migrating from Docker to Podman can be found at
#	https://podman-desktop.io/docs/migrating-from-docker.
#
#	Run these commands to install everything needed.
#	$ make dev-brew
#	$ make dev-docker
#	$ make dev-gotooling

# ==============================================================================
# Running Test
#
#	Running the tests is a good way to verify you have installed most of the
#	dependencies properly.
#
#	$ make test

# ==============================================================================
# Running The Project
#
#	$ make dev-up
#	$ make dev-update-apply
#	$ make token
#	$ export TOKEN=<token>
#	$ make users
#
#	You can use `make dev-status` to look at the status of your KIND cluster.

# ==============================================================================
# CLASS NOTES
#
# Kind
# 	For full Kind v0.31 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.31.0
#
# RSA Keys
# 	To generate a private/public key PEM file.
# 	$ openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# 	$ openssl rsa -pubout -in private.pem -out public.pem
# 	$ ./admin genkey
#
# Testing Coverage
# 	$ go test -coverprofile p.out
# 	$ go tool cover -html p.out
#
# Module Call Examples
# 	$ curl https://proxy.golang.org/github.com/ardanlabs/conf/@v/list
# 	$ curl https://proxy.golang.org/github.com/ardanlabs/conf/v3/@v/list
# 	$ curl https://proxy.golang.org/github.com/ardanlabs/conf/v3/@v/v3.1.1.info
# 	$ curl https://proxy.golang.org/github.com/ardanlabs/conf/v3/@v/v3.1.1.mod
# 	$ curl https://proxy.golang.org/github.com/ardanlabs/conf/v3/@v/v3.1.1.zip
# 	$ curl https://sum.golang.org/lookup/github.com/ardanlabs/conf/v3@v3.1.1
#
# OPA Playground
# 	https://play.openpolicyagent.org/
# 	https://academy.styra.com/
# 	https://www.openpolicyagent.org/docs/latest/policy-reference/

# ==============================================================================
# Define dependencies

GOLANG          := golang:1.26
ALPINE          := alpine:3.23
KIND            := kindest/node:v1.35.0
POSTGRES        := postgres:18.3
GRAFANA         := grafana/grafana:12.4.0
PROMETHEUS      := prom/prometheus:v3.10.0
TEMPO           := grafana/tempo:2.10.0
LOKI            := grafana/loki:3.6.0
PROMTAIL        := grafana/promtail:3.6.0

KIND_CLUSTER    := schoolcrm-cluster
NAMESPACE       := schoolcrm-system
KUBE_CONTEXT    ?= kind-$(KIND_CLUSTER)
KUBECTL         ?= kubectl --context $(KUBE_CONTEXT)

SCHOOLCRM_APP   := schoolcrm
AUTH_APP        := auth
WEB_ADMIN_APP   := web-admin
RAG_APP         := rag

BASE_IMAGE_NAME := localhost/schoolcrm

VERSION         := 0.0.1

SCHOOLCRM_IMAGE := $(BASE_IMAGE_NAME)/$(SCHOOLCRM_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)
WEB_ADMIN_IMAGE := $(BASE_IMAGE_NAME)/$(WEB_ADMIN_APP):$(VERSION)
RAG_IMAGE       := $(BASE_IMAGE_NAME)/$(RAG_APP):$(VERSION)

# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"

# ==============================================================================
# Detect operating system and set the appropriate open command

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	OPEN_CMD := open
else
	OPEN_CMD := xdg-open
endif

# ==============================================================================
# Install dependencies

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

dev-brew:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list watch || brew install watch
	brew list protobuf || brew install protobuf
	brew list grpcurl || brew install grpcurl

dev-docker:
	docker pull docker.io/$(GOLANG) & \
	docker pull docker.io/$(ALPINE) & \
	docker pull docker.io/$(KIND) & \
	docker pull docker.io/$(POSTGRES) & \
	docker pull docker.io/$(GRAFANA) & \
	docker pull docker.io/$(PROMETHEUS) & \
	docker pull docker.io/$(TEMPO) & \
	docker pull docker.io/$(LOKI) & \
	docker pull docker.io/$(PROMTAIL) & \
	wait;

# ==============================================================================
# Building containers

build: schoolcrm metrics auth rag web-admin

schoolcrm:
	docker build \
		-f zarf/docker/dockerfile.schoolcrm \
		-t $(SCHOOLCRM_IMAGE) \
		--build-arg BUILD_TAG=$(VERSION) \
		--build-arg BUILD_DATE=$$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

metrics:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t $(METRICS_IMAGE) \
		--build-arg BUILD_TAG=$(VERSION) \
		--build-arg BUILD_DATE=$$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

auth:
	docker build \
		-f zarf/docker/dockerfile.auth \
		-t $(AUTH_IMAGE) \
		--build-arg BUILD_TAG=$(VERSION) \
		--build-arg BUILD_DATE=$$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

rag:
	docker build \
		-f zarf/docker/dockerfile.rag \
		-t $(RAG_IMAGE) \
		--build-arg BUILD_TAG=$(VERSION) \
		--build-arg BUILD_DATE=$$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

web-admin:
	docker build \
		-f zarf/docker/dockerfile.web-admin \
		-t $(WEB_ADMIN_IMAGE) \
		--build-arg BUILD_TAG=$(VERSION) \
		--build-arg BUILD_DATE=$$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==============================================================================
# Running from within k8s/kind
# Docker Desktop 28.3.2 changed how it stores image layers, causing KIND's kind
# load docker-image command to fail with "content digest not found" errors. The
# workaround uses docker save | docker exec to bypass this incompatibility for
# the critical images allowing this to work without a network.

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	$(KUBECTL) wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	docker save $(POSTGRES) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
	docker save $(GRAFANA) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
	docker save $(PROMETHEUS) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
	docker save $(TEMPO) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
	docker save $(LOKI) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
	docker save $(PROMTAIL) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import - & \
	wait;

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-up-registry:
	KIND_CLUSTER_NAME=$(KIND_CLUSTER) ./zarf/k8s/dev/kind-with-registry.sh

dev-down-registry:
	kind delete cluster --name $(KIND_CLUSTER)
	docker stop kind-registry 2>/dev/null || true
	docker rm kind-registry 2>/dev/null || true

dev-status-all:
	$(KUBECTL) get nodes -o wide
	$(KUBECTL) get svc -o wide
	$(KUBECTL) get pods -o wide --watch --all-namespaces

dev-status:
	watch -n 2 $(KUBECTL) get pods -o wide --all-namespaces

# ------------------------------------------------------------------------------

dev-load:
	kind load docker-image $(SCHOOLCRM_IMAGE) --name $(KIND_CLUSTER) & \
	kind load docker-image $(METRICS_IMAGE) --name $(KIND_CLUSTER) & \
	kind load docker-image $(AUTH_IMAGE) --name $(KIND_CLUSTER) & \
	kind load docker-image $(WEB_ADMIN_IMAGE) --name $(KIND_CLUSTER) & \
	kind load docker-image $(RAG_IMAGE) --name $(KIND_CLUSTER) & \
	wait;

dev-apply: dev-load
	kustomize build zarf/k8s/dev/grafana | $(KUBECTL) apply -f -
	kustomize build zarf/k8s/dev/prometheus | $(KUBECTL) apply -f -
	kustomize build zarf/k8s/dev/tempo | $(KUBECTL) apply -f -
	kustomize build zarf/k8s/dev/loki | $(KUBECTL) apply -f -
	kustomize build zarf/k8s/dev/promtail | $(KUBECTL) apply -f -

	kustomize build zarf/k8s/dev/database | $(KUBECTL) apply -f -
	$(KUBECTL) rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	kustomize build zarf/k8s/dev/auth | $(KUBECTL) apply -f -
	$(KUBECTL) wait pods --namespace=$(NAMESPACE) --selector app=$(AUTH_APP) --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/rag | $(KUBECTL) apply -f -
	$(KUBECTL) wait pods --namespace=$(NAMESPACE) --selector app=$(RAG_APP) --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/schoolcrm | $(KUBECTL) apply -f -
	$(KUBECTL) wait pods --namespace=$(NAMESPACE) --selector app=$(SCHOOLCRM_APP) --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/web-admin | $(KUBECTL) apply -f -
	$(KUBECTL) wait pods --namespace=$(NAMESPACE) --selector app=$(WEB_ADMIN_APP) --timeout=120s --for=condition=Ready

dev-restart:
	$(KUBECTL) rollout restart deployment $(AUTH_APP) --namespace=$(NAMESPACE)
	$(KUBECTL) rollout restart deployment $(RAG_APP) --namespace=$(NAMESPACE)
	$(KUBECTL) rollout restart deployment $(SCHOOLCRM_APP) --namespace=$(NAMESPACE)
	$(KUBECTL) rollout restart deployment $(WEB_ADMIN_APP) --namespace=$(NAMESPACE)

dev-run: build dev-up dev-load dev-apply

dev-update: build dev-load dev-restart

dev-update-apply: build dev-load dev-apply

dev-logs:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=$(SCHOOLCRM_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run api/tooling/logfmt/main.go -service=$(SCHOOLCRM_APP)

dev-logs-auth:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=$(AUTH_APP) --all-containers=true -f --tail=100 | go run api/tooling/logfmt/main.go

dev-logs-web-admin:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=$(WEB_ADMIN_APP) --all-containers=true -f --tail=100

dev-logs-rag:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=$(RAG_APP) --all-containers=true -f --tail=100

# ------------------------------------------------------------------------------

dev-logs-init:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=$(SCHOOLCRM_APP) -f --tail=100 -c init-migrate-seed

dev-describe-node:
	$(KUBECTL) describe node

dev-describe-deployment:
	$(KUBECTL) describe deployment --namespace=$(NAMESPACE) $(SCHOOLCRM_APP)

dev-describe-schoolcrm:
	$(KUBECTL) describe pod --namespace=$(NAMESPACE) -l app=$(SCHOOLCRM_APP)

dev-describe-auth:
	$(KUBECTL) describe pod --namespace=$(NAMESPACE) -l app=$(AUTH_APP)

dev-describe-database:
	$(KUBECTL) describe pod --namespace=$(NAMESPACE) -l app=database

dev-describe-grafana:
	$(KUBECTL) describe pod --namespace=$(NAMESPACE) -l app=grafana

# ------------------------------------------------------------------------------

dev-logs-db:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=database --all-containers=true -f --tail=100

dev-logs-grafana:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=grafana --all-containers=true -f --tail=100

dev-logs-tempo:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=tempo --all-containers=true -f --tail=100

dev-logs-loki:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=loki --all-containers=true -f --tail=100

dev-logs-promtail:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=promtail --all-containers=true -f --tail=100

# ------------------------------------------------------------------------------

dev-services-delete:
	kustomize build zarf/k8s/dev/rag | $(KUBECTL) delete -f -
	kustomize build zarf/k8s/dev/schoolcrm | $(KUBECTL) delete -f -
	kustomize build zarf/k8s/dev/grafana | $(KUBECTL) delete -f -
	kustomize build zarf/k8s/dev/tempo | $(KUBECTL) delete -f -
	kustomize build zarf/k8s/dev/loki | $(KUBECTL) delete -f -
	kustomize build zarf/k8s/dev/promtail | $(KUBECTL) delete -f -
	kustomize build zarf/k8s/dev/database | $(KUBECTL) delete -f -

dev-describe-replicaset:
	$(KUBECTL) get rs
	$(KUBECTL) describe rs --namespace=$(NAMESPACE) -l app=$(SCHOOLCRM_APP)

dev-events:
	$(KUBECTL) get ev --sort-by metadata.creationTimestamp

dev-events-warn:
	$(KUBECTL) get ev --field-selector type=Warning --sort-by metadata.creationTimestamp

dev-shell:
	$(KUBECTL) exec --namespace=$(NAMESPACE) -it $(shell $(KUBECTL) get pods --namespace=$(NAMESPACE) | grep schoolcrm | cut -c1-26) --container $(SCHOOLCRM_APP) -- /bin/sh

dev-auth-shell:
	$(KUBECTL) exec --namespace=$(NAMESPACE) -it $(shell $(KUBECTL) get pods --namespace=$(NAMESPACE) | grep auth | cut -c1-26) --container $(AUTH_APP) -- /bin/sh

dev-db-shell:
	$(KUBECTL) exec --namespace=$(NAMESPACE) -it $(shell $(KUBECTL) get pods --namespace=$(NAMESPACE) | grep database | cut -c1-10) -- /bin/sh

dev-database-restart:
	$(KUBECTL) rollout restart statefulset database --namespace=$(NAMESPACE)

# ==============================================================================
# Docker Compose

compose-up:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml -p compose up -d

compose-build-up: build compose-up

compose-down:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml down

compose-logs:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml logs

# ==============================================================================
# Running locally without Docker or Kubernetes

local-db-url:
	@echo "PostgreSQL must be running on localhost:5454 with user/password postgres/postgres."
	@echo "If you do want Docker for only the database: docker compose -f zarf/compose/docker_compose.yaml -p compose up -d database"

local-migrate:
	SCHOOLCRM_DB_HOST=localhost:5454 \
	SCHOOLCRM_DB_DISABLE_TLS=true \
	go run api/tooling/admin/main.go migrate

local-seed: local-migrate
	SCHOOLCRM_DB_HOST=localhost:5454 \
	SCHOOLCRM_DB_DISABLE_TLS=true \
	go run api/tooling/admin/main.go seed

local-auth:
	AUTH_DB_HOST=localhost:5454 \
	AUTH_DB_DISABLE_TLS=true \
	AUTH_WEB_API_HOST=0.0.0.0:6000 \
	AUTH_WEB_GRPC_HOST=0.0.0.0:6001 \
	AUTH_WEB_DEBUG_HOST=0.0.0.0:6010 \
	go run api/services/auth/main.go | go run api/tooling/logfmt/main.go

local-schoolcrm:
	SCHOOLCRM_DB_HOST=localhost:5454 \
	SCHOOLCRM_DB_DISABLE_TLS=true \
	SCHOOLCRM_AUTH_HOST=http://localhost:6000 \
	SCHOOLCRM_AUTH_GRPC_HOST=localhost:6001 \
	SCHOOLCRM_WEB_API_HOST=0.0.0.0:3000 \
	SCHOOLCRM_WEB_DEBUG_HOST=0.0.0.0:3010 \
	go run api/services/schoolcrm/main.go | go run api/tooling/logfmt/main.go

local-rag:
	cd api/services/RAG && \
	uv sync --extra dev --locked && \
	RAG_API_HOST=0.0.0.0 \
	RAG_API_PORT=4545 \
	RAG_AUTH_SERVICE_URL=http://localhost:6000 \
	RAG_ALLOW_ANONYMOUS=false \
	RAG_FILE_STORAGE_DIR=./var/files \
	RAG_GRAPH_RETRIEVER=neo4j \
	RAG_NEO4J_URI=bolt://localhost:7687 \
	RAG_NEO4J_USERNAME=neo4j \
	RAG_NEO4J_PASSWORD=password \
	RAG_ADMISSIONS_ANSWER_PROVIDER=ollama \
	RAG_OLLAMA_BASE_URL=http://localhost:11434 \
	RAG_OLLAMA_MODEL=nemotron-3-super:cloud \
	uv run uvicorn main:app --reload --host 0.0.0.0 --port 4545

local-rag-dev:
	cd api/services/RAG && \
	uv sync --extra dev --locked && \
	RAG_API_HOST=0.0.0.0 \
	RAG_API_PORT=4545 \
	RAG_ALLOW_ANONYMOUS=true \
	RAG_FILE_STORAGE_DIR=./var/files \
	RAG_GRAPH_RETRIEVER=memory \
	RAG_ADMISSIONS_ANSWER_PROVIDER=stub \
	uv run uvicorn main:app --reload --host 0.0.0.0 --port 4545

local-rag-check:
	cd api/services/RAG && \
	uv sync --extra dev --locked && \
	uv run ruff check . && \
	uv run python -m compileall -q adapters domain infrastructure use_cases main.py && \
	uv run pytest tests -v

local-web-admin:
	npm --prefix api/frontends/web-admin start

local-run-help:
	@echo "Run these targets in separate terminals:"
	@echo "  make local-db-url       # explains localhost database requirement"
	@echo "  make local-seed         # migrate and seed localhost database"
	@echo "  make local-auth         # auth API on :6000, gRPC on :6001"
	@echo "  make local-schoolcrm    # SchoolCRM API on :3000"
	@echo "  make local-rag          # RAG API on :4545 using auth service"
	@echo "  make local-rag-dev      # RAG API on :4545 with anonymous auth"
	@echo "  make local-rag-check    # validate RAG lint, compile, and tests"
	@echo "  make local-web-admin    # Angular admin app on :4200"

# ==============================================================================
# Administration

migrate:
	export SCHOOLCRM_DB_HOST=localhost; go run api/tooling/admin/main.go migrate

seed: migrate
	export SCHOOLCRM_DB_HOST=localhost; go run api/tooling/admin/main.go seed

pgcli:
	pgcli postgresql://postgres:postgres@localhost

liveness:
	curl -i http://localhost:3000/v1/liveness

readiness:
	curl -i http://localhost:3000/v1/readiness

token-gen:
	export SCHOOLCRM_DB_HOST=localhost; go run api/tooling/admin/main.go gentoken 5cf37266-3473-4006-984f-9325122678b7 54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

# ==============================================================================
# Metrics and Tracing

metrics-view-sc:
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

metrics-view:
	expvarmon -ports="localhost:4020" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

grafana:
	$(OPEN_CMD) http://localhost:3100/

statsviz:
	$(OPEN_CMD) http://localhost:3010/debug/statsviz

# ==============================================================================
# Running tests within the local computer

test-down:
	docker stop servicetest
	docker rm servicetest -v

test-r:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test: test-only lint vuln-check

test-race: test-r lint vuln-check

# ==============================================================================
# Hitting endpoints

token:
	curl -s \
	--user "admin@example.com:gophers" http://localhost:6000/v1/auth/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

# export TOKEN="COPY TOKEN STRING FROM LAST CALL"

###
# Use this command to export the TOKEN variable directly into your shell
# To do so, you have to run: `eval $(make export-token)`
###
export-token:
	@$(MAKE) token \
	| sed -n 's/.*"token":"\([^"]*\)".*/export TOKEN="\1"/p'

token-grpc:
	grpcurl -plaintext -d '{"kid":"54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"}' \
      -H "Authorization: Basic YWRtaW5AZXhhbXBsZS5jb206Z29waGVycw==" \
      localhost:6001 auth.Auth/Token

users:
	curl -i \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/v1/users?page=1&rows=2"

users-timeout:
	curl -i \
	--max-time 1 \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/v1/users?page=1&rows=2"

load:
	hey -m GET -c 100 -n 1000 \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/v1/users?page=1&rows=2"

otel-test:
	curl -i \
	-H "Traceparent: 00-918dd5ecf264712262b68cf2ef8b5239-896d90f23f69f006-01" \
	--user "admin@example.com:gophers" http://localhost:6000/v1/auth/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

# ==============================================================================
# Protobuf support

authapp-proto-gen:
	protoc --go_out=app/domain/grpcauthapp --go_opt=paths=source_relative \
		--go-grpc_out=app/domain/grpcauthapp --go-grpc_opt=paths=source_relative \
		--proto_path=app/domain/grpcauthapp \
		app/domain/grpcauthapp/authapp.proto

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Class Stuff

run-auth:
	go run api/services/auth/main.go | go run api/tooling/logfmt/main.go

run:
	go run api/services/schoolcrm/main.go | go run api/tooling/logfmt/main.go

run-help:
	go run api/services/schoolcrm/main.go --help | go run api/tooling/logfmt/main.go

curl:
	curl -i http://localhost:3000/v1/hack

curl-auth:
	curl -i -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/hackauth

load-hack:
	hey -m GET -c 100 -n 100000 "http://localhost:3000/v1/hack"

admin:
	go run api/tooling/admin/main.go

ready:
	curl -i http://localhost:3000/v1/readiness

live:
	curl -i http://localhost:3000/v1/liveness

curl-create:
	curl -i -X POST \
	-H "Authorization: Bearer ${TOKEN}" \
	-H 'Content-Type: application/json' \
	-d '{"name":"bill","email":"b@gmail.com","roles":["ADMIN"],"department":"ITO","password":"123","passwordConfirm":"123"}' \
	http://localhost:3000/v1/users

# ==============================================================================
# Talk commands

talk-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	$(KUBECTL) wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	docker save $(POSTGRES) | docker exec -i $(KIND_CLUSTER)-control-plane ctr --namespace=k8s.io images import -

talk-load:
	kind load docker-image $(SCHOOLCRM_IMAGE) --name $(KIND_CLUSTER) & \
	kind load docker-image $(METRICS_IMAGE) --name $(KIND_CLUSTER) & \
	kind load docker-image $(AUTH_IMAGE) --name $(KIND_CLUSTER) & \
	wait;

talk-apply: talk-load
	kustomize build zarf/k8s/dev/database | $(KUBECTL) apply -f -
	$(KUBECTL) rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	kustomize build zarf/k8s/dev/auth | $(KUBECTL) apply -f -
	$(KUBECTL) wait pods --namespace=$(NAMESPACE) --selector app=$(AUTH_APP) --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/schoolcrm | $(KUBECTL) apply -f -
	$(KUBECTL) wait pods --namespace=$(NAMESPACE) --selector app=$(SCHOOLCRM_APP) --timeout=120s --for=condition=Ready

talk-run: build talk-up talk-load talk-apply

talk-run-load:
	hey -m GET -c 10 -n 1000 -H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/v1/users?page=1&rows=2"

talk-logs:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=$(SCHOOLCRM_APP) --all-containers=true -f --tail=100 --max-log-requests=6

talk-logs-cpu:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=$(SCHOOLCRM_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | grep SCHED

talk-logs-mem:
	$(KUBECTL) logs --namespace=$(NAMESPACE) -l app=$(SCHOOLCRM_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | grep "ms clock"

talk-describe:
	$(KUBECTL) describe pod --namespace=$(NAMESPACE) -l app=$(SCHOOLCRM_APP)

talk-describe-min:
	$(KUBECTL) describe pod --namespace=$(NAMESPACE) -l app=$(SCHOOLCRM_APP) | grep -i -A 21 '  schoolcrm:' | sed -n '1p;3p;11,16p;18,22p'

talk-metrics:
	expvarmon -ports="localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

# ==============================================================================
# Admin Frontend

ADMIN_FRONTEND_PREFIX := ./api/frontends/admin

write-token-to-env:
	echo "VITE_SERVICE_API=http://localhost:3000/v1" > ${ADMIN_FRONTEND_PREFIX}/.env
	make token | grep -o '"ey.*"' | awk '{print "VITE_SERVICE_TOKEN="$$1}' >> ${ADMIN_FRONTEND_PREFIX}/.env

admin-gui-install:
	pnpm -C ${ADMIN_FRONTEND_PREFIX} install

admin-gui-update:
	pnpm -C ${ADMIN_FRONTEND_PREFIX} update

admin-gui-dev: admin-gui-install
	pnpm -C ${ADMIN_FRONTEND_PREFIX} run dev

admin-gui-build: admin-gui-install
	pnpm -C ${ADMIN_FRONTEND_PREFIX} run build

admin-gui-start-build: admin-gui-build
	pnpm -C ${ADMIN_FRONTEND_PREFIX} run preview

admin-gui-run: write-token-to-env admin-gui-start-build

# ==============================================================================
# Help command
help:
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@echo "  dev-gotooling           Install Go tooling"
	@echo "  dev-brew                Install brew dependencies"
	@echo "  dev-docker              Pull Docker images"
	@echo "  build                   Build all the containers"
	@echo "  schoolcrm                   Build the SchoolCRM container"
	@echo "  metrics                 Build the metrics container"
	@echo "  auth                    Build the auth container"
	@echo "  rag                     Build the RAG container"
	@echo "  local-run-help          Show local non-Docker run commands"
	@echo "  local-auth              Run auth locally without Docker/Kubernetes"
	@echo "  local-schoolcrm         Run SchoolCRM API locally without Docker/Kubernetes"
	@echo "  local-rag               Run RAG locally against local auth"
	@echo "  local-rag-dev           Run RAG locally with anonymous auth"
	@echo "  local-web-admin         Run Angular admin locally"
	@echo "  dev-up                  Start the KIND cluster"
	@echo "  dev-down                Stop the KIND cluster"
	@echo "  dev-status-all          Show the status of the KIND cluster"
	@echo "  dev-status              Show the status of the pods"
	@echo "  dev-load                Load the containers into KIND"
	@echo "  dev-apply               Apply the manifests to KIND"
	@echo "  dev-restart             Restart the deployments"
	@echo "  dev-run              	 Build, up, load, and apply the deployments"
	@echo "  dev-update              Build, load, and restart the deployments"
	@echo "  dev-update-apply        Build, load, and apply the deployments"
	@echo "  dev-logs                Show the logs for the SchoolCRM service"
	@echo "  dev-logs-auth           Show the logs for the auth service"
	@echo "  dev-logs-rag            Show the logs for the RAG service"
	@echo "  dev-logs-init           Show the logs for the init container"
	@echo "  dev-describe-node       Show the node details"
	@echo "  dev-describe-deployment Show the deployment details"
	@echo "  dev-describe-schoolcrm      Show the SchoolCRM pod details"
	@echo "  dev-describe-auth       Show the auth pod details"
	@echo "  dev-describe-database   Show the database pod details"
	@echo "  dev-describe-grafana    Show the grafana pod details"
	@echo "  dev-logs-db             Show the logs for the database service"
	@echo "  dev-logs-grafana        Show the logs for the grafana service"
	@echo "  dev-logs-tempo          Show the logs for the tempo service"
	@echo "  dev-logs-loki           Show the logs for the loki service"
	@echo "  dev-logs-promtail       Show the logs for the promtail service"
	@echo "  dev-services-delete     Delete all"
	@echo "  dev-up-registry         Start KIND with local registry"
	@echo "  dev-down-registry       Stop KIND and local registry"
	@echo "  tilt-up                 Start Tilt"
	@echo "  tilt-up-hud             Start Tilt with HUD"
	@echo "  tilt-down               Stop Tilt"
	@echo "  tilt-clean              Stop Tilt and delete namespaces"
	@echo "  tilt-logs               Show logs for a Tilt service"
	@echo "  chart-install           Install a Helm chart"
	@echo "  chart-upgrade           Upgrade a Helm chart"
	@echo "  chart-uninstall         Uninstall a Helm chart"
	@echo "  chart-template          Render a Helm chart"
	@echo "  chart-lint              Lint a Helm chart"

# ==============================================================================
# ==============================================================================
# Development Workflow (Tilt)
# ==============================================================================

# Start local development with hot reload
tilt-up:
	tilt up

# Start Tilt with HUD
tilt-up-hud:
	tilt up --hud=true

# Stop Tilt
tilt-down:
	tilt down

# Clean up everything
tilt-clean:
	tilt down --delete-namespaces

# View Tilt logs for a specific service
tilt-logs:
	@echo "Usage: make tilt-logs SERVICE=schoolcrm"
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: SERVICE required"; \
		echo "Example: make tilt-logs SERVICE=schoolcrm"; \
		exit 1; \
	fi
	tilt logs $(SERVICE)

# ==============================================================================
# Individual Chart Deployment (Production)
# ==============================================================================

# Deploy individual chart
# Usage: make chart-install CHART=schoolcrm ENV=prod
chart-install:
	@if [ -z "$(CHART)" ]; then \
		echo "Error: CHART required"; \
		echo "Example: make chart-install CHART=schoolcrm ENV=prod"; \
		exit 1; \
	fi
	helm install $(CHART) ./zarf/helm/charts/$(CHART) \
		--namespace $(NAMESPACE) \
		--create-namespace \
		--values ./zarf/helm/charts/$(CHART)/values-$(ENV).yaml \
		--wait \
		--timeout 10m

# Upgrade individual chart
chart-upgrade:
	@if [ -z "$(CHART)" ]; then \
		echo "Error: CHART required"; \
		echo "Example: make chart-upgrade CHART=schoolcrm ENV=prod"; \
		exit 1; \
	fi
	helm upgrade $(CHART) ./zarf/helm/charts/$(CHART) \
		--namespace $(NAMESPACE) \
		--values ./zarf/helm/charts/$(CHART)/values-$(ENV).yaml \
		--wait \
		--timeout 10m

# Uninstall individual chart
chart-uninstall:
	@if [ -z "$(CHART)" ]; then \
		echo "Error: CHART required"; \
		echo "Example: make chart-uninstall CHART=schoolcrm"; \
		exit 1; \
	fi
	helm uninstall $(CHART) --namespace $(NAMESPACE) --ignore-not-found

# Template individual chart (debug)
chart-template:
	@if [ -z "$(CHART)" ]; then \
		echo "Error: CHART required"; \
		echo "Example: make chart-template CHART=schoolcrm ENV=dev"; \
		exit 1; \
	fi
	helm template $(CHART) ./zarf/helm/charts/$(CHART) \
		--values ./zarf/helm/charts/$(CHART)/values-$(ENV).yaml \
		--debug

# Lint individual chart
chart-lint:
	@if [ -z "$(CHART)" ]; then \
		echo "Error: CHART required"; \
		echo "Example: make chart-lint CHART=schoolcrm"; \
		exit 1; \
	fi
	helm lint ./zarf/helm/charts/$(CHART)
