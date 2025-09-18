FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o odc \
    ./cmd/odc

FROM ubuntu:22.04

# for NFS Client
RUN apt update && \
    DEBIAN_FRONTEND=noninteractive apt install -y \
    nfs-common \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/odc .

# Get docker binary to run docker inside the container
COPY --from=docker:cli /usr/local/bin/docker /usr/local/bin/docker

CMD ["/app/odc"]