#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

set -a
source /root/.secrets/ai-resume-analyzer
set +a

docker compose build
docker compose up -d

docker compose ps