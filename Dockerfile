FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
--mount=type=cache,target=/root/.cache/go-build \
go build -o /scansched ./cmd/scansched


FROM debian:stable-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /data
COPY --from=builder /scansched /usr/local/bin/scansched
COPY configs/sample.yaml /etc/scansched.yaml
ENTRYPOINT ["/usr/local/bin/scansched","--config","/etc/scansched.yaml"]