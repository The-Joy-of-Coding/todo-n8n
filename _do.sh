#!/bin/bash

info() {
  echo -e "\033[1;34m[INFO]\033[0m $1";
}

err() {
  echo -e "\033[1;31m[ERROR]\033[0m $1"; exit 1;
}

export INSECURE_API="http://localhost"
export SECURE_API="https://example.com"
export BINARY_NAME="todo-n8n"

case "$1" in
  test)
    info "Running tests..."
    shift
    go test ./... "$@" || err "Test Failed!"
    ;;
  build)
    info "Building binary: $BINARY_NAME"
    go build -ldflags="-s -w" -o "$BINARY_NAME" . || err "Build failed."
    SIZE_BEFORE=$(du -m "$BINARY_NAME" | cut -f1)
    info "Compressing binary..."
    upx -1 "$BINARY_NAME" > /dev/null 2>&1 || err "Compression failed."
    SIZE_AFTER=$(du -m "$BINARY_NAME" | cut -f1)
    echo -e "\033[1;32m[DONE]\033[0m Size reduced: ${SIZE_BEFORE}MB -> ${SIZE_AFTER}MB"
    ;;
  rm)
    info "Removing binary: $BINARY_NAME"
    rm -f "$BINARY_NAME"
    ;;
  run)
    shift
    go run . "$@"
    ;;
  *)
    echo "Usage: $0 {test|build|run} [options]"
    echo ""
    echo "Commands:"
    echo "  test   Run all tests (pass -v for verbose)"
    echo "  build  Compile the Go binary"
    echo "  run    Execute the app (pass -add, -list, etc.)"
    echo "  rm     Removes the built binary"
    exit 1
esac
