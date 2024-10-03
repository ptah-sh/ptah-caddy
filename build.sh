#!/usr/bin/env sh

set -e

cd observability

go mod tidy

go test

cd ..

go mod tidy

# TODO: do not use master
xcaddy build master --with github.com/ptah-sh/ptah-caddy/observability=./observability
