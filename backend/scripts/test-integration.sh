#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
container="vinhomes-payment-tests-${RANDOM}"
trap 'docker rm -f -v "$container" >/dev/null 2>&1 || true' EXIT
docker run -d --rm --name "$container" -e POSTGRES_PASSWORD=test-only -e POSTGRES_DB=vinhomes_test -p 127.0.0.1::5432 postgres:16-alpine >/dev/null
for _ in {1..60}; do
  if docker exec "$container" pg_isready -U postgres -d vinhomes_test >/dev/null 2>&1; then break; fi
  sleep 0.25
done
port=$(docker port "$container" 5432/tcp)
export TEST_DATABASE_URL="postgres://postgres:test-only@${port}/vinhomes_test?sslmode=disable"
go test -race -count=1 ./...
