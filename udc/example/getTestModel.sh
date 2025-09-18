#!/bin/bash

set -e

cd "$(dirname "$0")"

docker pull busybox
docker save busybox > model.tar
docker image rm busybox