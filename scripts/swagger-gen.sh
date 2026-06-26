#!/usr/bin/env bash
# Generate embedded Swagger assets from the canonical contract (docs/contracts/openapi.yaml).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

OPENAPI_SRC="${OPENAPI_SRC:-docs/contracts/openapi.yaml}"
OPENAPI_EMBED="${OPENAPI_EMBED:-internal/api/swagger/openapi.yaml}"
ADMIN_SPEC="${ADMIN_SPEC:-internal/api/swagger/admin-openapi.yaml}"
GO="${GO:-go}"

die() { echo "swagger-gen: $*" >&2; exit 1; }

[[ -f "$OPENAPI_SRC" ]] || die "contract not found: $OPENAPI_SRC"

validate_yaml() {
  local f="$1"
  if command -v yq >/dev/null 2>&1; then
    yq eval '.' "$f" >/dev/null
  elif command -v python3 >/dev/null 2>&1; then
    python3 -c "import yaml; yaml.safe_load(open('$f'))"
  else
    die "install yq or python3+pyyaml to validate YAML"
  fi
}

echo "==> validate YAML syntax"
validate_yaml "$OPENAPI_SRC"
[[ -f "$ADMIN_SPEC" ]] && validate_yaml "$ADMIN_SPEC"

echo "==> validate OpenAPI (kin-openapi parse + structure)"
$GO run ./tools/openapi-validate "$OPENAPI_SRC"
[[ -f "$ADMIN_SPEC" ]] && $GO run ./tools/openapi-validate "$ADMIN_SPEC"

if [[ "${OPENAPI_VALIDATE_STRICT:-}" == "1" ]]; then
  echo "==> strict OpenAPI validation (OPENAPI_VALIDATE_STRICT=1)"
  $GO run ./tools/openapi-validate --strict "$OPENAPI_SRC"
  [[ -f "$ADMIN_SPEC" ]] && $GO run ./tools/openapi-validate --strict "$ADMIN_SPEC"
fi

echo "==> embed spec for Swagger UI"
mkdir -p "$(dirname "$OPENAPI_EMBED")"
cp "$OPENAPI_SRC" "$OPENAPI_EMBED"

if ! diff -q "$OPENAPI_SRC" "$OPENAPI_EMBED" >/dev/null; then
  die "embed mismatch after copy"
fi

echo "swagger-gen: OK ($OPENAPI_SRC -> $OPENAPI_EMBED)"
