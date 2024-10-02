#!/usr/bin/env bash

set -e

./build.sh

./caddy run --config Caddyfile.dev --adapter caddyfile
