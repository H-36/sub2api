#!/bin/sh
set -eu

SCRIPT_DIR="$(CDPATH= cd "$(dirname "$0")" && pwd -P)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
APP_DIR="${REPO_ROOT}/third_party/gpt_image_playground"
FRONTEND_DIR="${REPO_ROOT}/frontend"
OUTPUT_DIR="${FRONTEND_DIR}/public/image-playground-app"

if [ ! -f "${APP_DIR}/package.json" ]; then
  echo "gpt_image_playground source not found at ${APP_DIR}" >&2
  exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
  echo "npm is required to build gpt_image_playground" >&2
  exit 1
fi

cd "${APP_DIR}"

if [ ! -d node_modules ]; then
  npm ci
fi

npm run build

rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}"
cp -a "${APP_DIR}/dist/." "${OUTPUT_DIR}/"

# Safety cleanup for older local checkouts.
rm -f "${OUTPUT_DIR}/sw.js" "${OUTPUT_DIR}/manifest.webmanifest"

echo "Image playground built to ${OUTPUT_DIR}"
