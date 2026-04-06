#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-${SCRIPT_DIR}/docker-compose.local.yml}"
ENV_FILE="${ENV_FILE:-${SCRIPT_DIR}/.env}"

echo "Using image: $(grep -E '^SUB2API_IMAGE=' "${ENV_FILE}" | cut -d= -f2- || echo 'not set')"
docker compose --env-file "${ENV_FILE}" -f "${COMPOSE_FILE}" pull sub2api
docker compose --env-file "${ENV_FILE}" -f "${COMPOSE_FILE}" up -d --no-deps sub2api

echo "Sub2API updated."
