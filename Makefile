APP_NAME := markdown-sync
BINDIR := bin
DISTDIR := dist

.PHONY: all build test fmt vet lint tidy setup docker-image cross clean

all: build

build:
	@mkdir -p $(BINDIR)
	go build -v -o $(BINDIR)/$(APP_NAME) ./cmd/markdown-sync

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
	docker build -t $(APP_NAME):dev .

cross:
	@mkdir -p $(DISTDIR)
	GOOS=linux GOARCH=amd64 go build -v -o $(DISTDIR)/$(APP_NAME)-linux-amd64 ./cmd/markdown-sync

clean:
	rm -rf $(BINDIR) $(DISTDIR)
