FROM golang:1.26-alpine AS builder

WORKDIR /src

ENV CGO_ENABLED=0 \
    GOFLAGS=-mod=readonly

# Cache deps layer
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -ldflags="-s -w" -o /out/iustitia ./cmd/app

RUN mkdir -p /out/data /out/app-secrets

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/iustitia                     /app/iustitia
COPY --from=builder /src/migrations                   /app/migrations
COPY --from=builder --chown=nonroot:nonroot /out/data /data

ENV HTTP_ADDR=:8080 \
    DB_PATH=/data/iustitia.db \
    MIGRATIONS_DIR=/app/migrations \
    SECRETS_DIR=/app/secrets \
    ENV=prod

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/iustitia"]
