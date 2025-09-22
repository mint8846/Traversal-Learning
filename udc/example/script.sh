#!/bin/bash
set -e
echo "$1" >> /tmp/result/execute_log
# Pass file content via stdin. In case of docker-in-docker, SRC file mount path can only use actual host paths.
# Therefore, since there are constraints on mounting files inside containers, file content is passed via stdin
cat "$1"|docker run --rm --name "test-model" -i busybox cat >> /tmp/result/execute_log
rm -rf "$1"
echo "======" >> /tmp/result/execute_log
