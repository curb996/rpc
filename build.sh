#!/usr/bin/env bash
set -e
mode=\$1  # so / exe
case $mode in
  so)
    GOOS=\$2 GOARCH=\$3 go build -buildmode=plugin -o ../../../../plugins/modbus.so
    ;;
  exe)
    GOOS=\$2 GOARCH=\$3 go build -o ../../../../plugins/modbus_plugin.exe ./cmd
    ;;
esac
