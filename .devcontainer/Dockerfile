FROM golang:1.21-bullseye

RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    make \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install staticcheck
RUN curl -sSfL https://raw.githubusercontent.com/dominikh/go-tools-install-script/master/install.sh | bash -s -- -b /usr/local/bin staticcheck || true

# Create a non-root user for VS Code remote
RUN useradd -m vscode || true
USER vscode
WORKDIR /workspace
