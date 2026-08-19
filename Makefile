APP_NAME := gdocs-markdown-sync
BINDIR := bin
DISTDIR := dist
DOCKER_IMAGE ?= $(APP_NAME):dev
DOCKERFILE := .devcontainer/Dockerfile
DOCKER_CONTEXT := .devcontainer

IN_CONTAINER := $(shell if [ -f /.dockerenv ] || [ "$${DEVCONTAINER}" = "true" ]; then echo 1; fi)
HOST_OK := docker-image docker-run clean

ifeq ($(IN_CONTAINER),)
  ifeq ($(filter $(HOST_OK),$(MAKECMDGOALS)),)
    $(error This must run inside the .devcontainer image. Use: make docker-run CMD='make $(or $(MAKECMDGOALS),all)'. Do not install Go on the host. See AGENTS.md.)
  endif
endif

.PHONY: all build test fmt vet lint tidy setup docker-image docker-run cross clean

all: build

build:
	@mkdir -p $(BINDIR)
	go build -v -o $(BINDIR)/$(APP_NAME) ./cmd/gdocs-markdown-sync

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

tidy:
	go mod tidy

lint:
	if command -v staticcheck >/dev/null 2>&1; then staticcheck ./...; else echo "staticcheck not found; run 'make setup'"; fi

setup:
	@echo "Installing developer tools..."
	@echo "Installing staticcheck (if possible)"
	go install honnef.co/go/tools/cmd/staticcheck@latest || true

docker-image:
	docker build -t $(DOCKER_IMAGE) -f $(DOCKERFILE) $(DOCKER_CONTEXT)

# Usage (from the host): make docker-run CMD='make test'
docker-run: docker-image
	@test -n "$(CMD)" || (echo "CMD is required, e.g. make docker-run CMD='make test'" >&2; exit 1)
	docker run --rm \
		-v "$(CURDIR):/workspace" \
		-w /workspace \
		-e GOCACHE=/tmp/go-cache \
		-e GOMODCACHE=/tmp/gomod \
		-e DEVCONTAINER=true \
		$(DOCKER_IMAGE) \
		$(CMD)

cross:
	@mkdir -p $(DISTDIR)
	GOOS=linux GOARCH=amd64 go build -v -o $(DISTDIR)/$(APP_NAME)-linux-amd64 ./cmd/gdocs-markdown-sync

clean:
	rm -rf $(BINDIR) $(DISTDIR)
