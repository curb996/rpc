#!/usr/bin/env bash
set -e
# build core
for GOARCH in amd64 arm64; do
  env GOOS=linux  GOARCH=$GOARCH CGO_ENABLED=0 \
    go build -o dist/collector_linux_${GOARCH} ./cmd/edge-collector
done
# build plugins
pushd internal/drivers/modbus
./build.sh so linux amd64
./build.sh so linux arm64
./build.sh exe windows amd64
popd
