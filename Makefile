APP_NAME := gdocs-markdown-sync
BINDIR := bin
DISTDIR := dist
DOCKER_IMAGE ?= $(APP_NAME):dev
DOCKERFILE := .devcontainer/Dockerfile
DOCKER_CONTEXT := .devcontainer

# True when already running inside the devcontainer (avoid docker-in-docker).
IN_CONTAINER := $(shell if [ -f /.dockerenv ] || [ "$${DEVCONTAINER}" = "true" ]; then echo 1; fi)

# On the host, run the command in the .devcontainer image. Inside the container, run it natively.
run_in_env = $(if $(IN_CONTAINER),$(1),$(MAKE) docker-run CMD="$(1)")

.PHONY: all build test fmt vet lint tidy setup docker-image docker-run cross clean

all: build

build:
	@mkdir -p $(BINDIR)
	$(call run_in_env,go build -v -o $(BINDIR)/$(APP_NAME) ./cmd/gdocs-markdown-sync)

test:
	$(call run_in_env,go test ./...)

fmt:
	$(call run_in_env,gofmt -w .)

vet:
	$(call run_in_env,go vet ./...)

tidy:
	$(call run_in_env,go mod tidy)

lint:
ifeq ($(IN_CONTAINER),1)
	if command -v staticcheck >/dev/null 2>&1; then staticcheck ./...; else echo "staticcheck not found; run 'make setup'"; fi
else
	$(MAKE) docker-run CMD='staticcheck ./...'
endif

setup:
	@echo "Installing developer tools..."
	@echo "Installing staticcheck (if possible)"
	$(call run_in_env,go install honnef.co/go/tools/cmd/staticcheck@latest)

docker-image:
	docker build -t $(DOCKER_IMAGE) -f $(DOCKERFILE) $(DOCKER_CONTEXT)

# Usage: make docker-run CMD='go test ./...'
docker-run: docker-image
	@test -n "$(CMD)" || (echo "CMD is required, e.g. make docker-run CMD='go test ./...'" >&2; exit 1)
	docker run --rm \
		-v "$(CURDIR):/workspace" \
		-w /workspace \
		-e GOCACHE=/tmp/go-cache \
		-e GOMODCACHE=/tmp/gomod \
		$(DOCKER_IMAGE) \
		$(CMD)

cross:
	@mkdir -p $(DISTDIR)
	$(call run_in_env,env GOOS=linux GOARCH=amd64 go build -v -o $(DISTDIR)/$(APP_NAME)-linux-amd64 ./cmd/gdocs-markdown-sync)

clean:
	rm -rf $(BINDIR) $(DISTDIR)
