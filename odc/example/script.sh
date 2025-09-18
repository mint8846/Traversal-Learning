#!/bin/bash
set -e
# Refer to env MODEL_OUTPUT_FILE_PATH
docker run --rm --name "test-model" busybox ls > /tmp/output
