#!/usr/bin/env sh

set -e

cd observability

go mod tidy

go test

cd ..

go mod tidy

xcaddy build --with github.com/ptah-sh/ptah-caddy/observability=./observability
