FROM caddy:2.9.1-builder AS builder

WORKDIR /app

COPY . .

RUN ./build.sh

FROM caddy:2.9.1

COPY --from=builder /app/caddy /usr/bin/caddy
