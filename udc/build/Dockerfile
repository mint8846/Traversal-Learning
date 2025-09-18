FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o udc \
    ./cmd/udc

FROM ubuntu:22.04

# for NFS
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    nfs-kernel-server \
    rpcbind \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /run/sendsigs.omit.d \
    && mkdir -p /var/lib/nfs/rpc_pipefs \
    && mkdir -p /var/lib/nfs/v4recovery \
    && mkdir -p /nfs/share

RUN echo 'STATDOPTS="--port 32765 --outgoing-port 32766"' >> /etc/default/nfs-common && \
    echo 'RPCMOUNTDOPTS="--manage-gids --port 20048"' >> /etc/default/nfs-kernel-server

ENV NFS_EXPORT_PATH=/nfs/share

COPY build/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /app

COPY --from=builder /app/udc .

# Get docker binary to run docker inside the container
COPY --from=docker:cli /usr/local/bin/docker /usr/local/bin/docker

EXPOSE 8080 111 2049 20048 32765 32766

#VOLUME ["/nfs/share"]

CMD ["/entrypoint.sh"]