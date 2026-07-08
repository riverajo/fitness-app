#!/bin/bash
set -e

echo "=== Installing golangci-lint ==="
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.63.4

echo "=== Installing Antigravity agy CLI ==="
curl -fsSL https://antigravity.google/cli/install.sh | bash

echo "=== Installing Forgejo CLI (fj) ==="
mkdir -p /tmp/fj-install
curl -sSL -o /tmp/fj-install/fj.tar.gz https://codeberg.org/forgejo-contrib/forgejo-cli/releases/download/v0.5.0/forgejo-cli-x86_64-linux.tar.gz
tar -xzf /tmp/fj-install/fj.tar.gz -C /tmp/fj-install
mv /tmp/fj-install/fj /usr/local/bin/fj
rm -rf /tmp/fj-install

echo "=== Setup complete ==="
