FROM caddy:2.8.4-builder AS builder

WORKDIR /app

COPY . .

RUN ./build.sh

FROM caddy:2.8.4

COPY --from=builder /app/caddy /usr/bin/caddy
