#!/usr/bin/env bash
#
# Generate Go code from the OpenAPI spec.
#   1. bundle split spec (openapi.yml + routes/ + components/) into single yaml
#   2. run oapi-codegen ×3 (types / server / spec)
#   3. cleanup bundle on exit (trap)
#
# Used by `make openapi` and CI. Works on linux/macos/windows (git-bash).

set -euo pipefail

# Change to backend root regardless of where the script was invoked from.
cd "$(dirname "$0")/.."

REDOCLY_VERSION="${REDOCLY_VERSION:-1.26.0}"
SPEC_SRC="internal/openapi/openapi.yml"
BUNDLE=$(mktemp -t iustitia-openapi-XXXXXX.yaml)

cleanup() {
    rm -f "$BUNDLE"
}
trap cleanup EXIT

# 1. Bundle split spec into single self-contained yaml.
command -v npx >/dev/null 2>&1 || {
    echo "error: npx (Node.js) required for redocly bundle. install Node.js or: npm install -g @redocly/cli" >&2
    exit 1
}

echo "-> bundling $SPEC_SRC via redocly v$REDOCLY_VERSION"
npx --yes "@redocly/cli@${REDOCLY_VERSION}" bundle "$SPEC_SRC" -o "$BUNDLE" --ext yaml

# 2. Run oapi-codegen for each artefact.
command -v oapi-codegen >/dev/null 2>&1 || {
    echo "error: oapi-codegen not found. install via: make install-tools" >&2
    exit 1
}

echo "-> generating types.gen.go"
oapi-codegen -config codegen/oapi-codegen-types.yml  "$BUNDLE"

echo "-> generating server.gen.go"
oapi-codegen -config codegen/oapi-codegen-server.yml "$BUNDLE"

echo "-> generating spec.gen.go"
oapi-codegen -config codegen/oapi-codegen-spec.yml   "$BUNDLE"

echo "✓ OpenAPI code generation completed"
