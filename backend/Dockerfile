# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.27-alpine AS build
WORKDIR /src

# Download dependencies first so this layer is cached until go.mod/go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# ---- Runtime stage ----
# Distroless: no shell or package manager, runs as non-root user.
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=build /out/server /app/server

ENV GIN_MODE=release \
    PORT=8080
EXPOSE 8080

USER nonroot:nonroot
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/app/server", "healthcheck"]
# Exec form so the server is PID 1 and receives SIGTERM for graceful shutdown.
ENTRYPOINT ["/app/server"]
