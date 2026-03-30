#!/usr/bin/env bash
# Builds the provider and runs the given terraform command inside test-local/.
# Usage:
#   ./run.sh plan
#   ./run.sh apply
#   ./run.sh destroy
#   ./run.sh show

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TEST_DIR="$REPO_ROOT/test-local"

echo "==> Building provider..."
(cd "$REPO_ROOT" && go build -o terraform-provider-netearthone .)
echo "    Binary: $REPO_ROOT/terraform-provider-netearthone"

echo "==> Running: terraform ${*:-plan}"
cd "$TEST_DIR"
TF_CLI_CONFIG_FILE="$TEST_DIR/.terraformrc" terraform "${@:-plan}"
