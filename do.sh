#!/bin/bash

info() {
  echo -e "\033[1;34m[INFO]\033[0m $1";
}

err() {
  echo -e "\033[1;31m[ERROR]\033[0m $1"; exit 1;
}

export INSECURE_API="http://localhost"
export SECURE_API="https://example.com"

case "$1" in
  test)
    info "Running tests..."
    shift
    go test ./... "$@" || err "Test Failed!";;
  build)
    info "Building binary: todo-n8n"
    go build -o todo-n8n . || err "Build failed.";;
  run)
    shift
    go run . "$@";;
  *)
    echo "Usage: $0 {test|build|run} [options]"
    echo ""
    echo "Commands:"
    echo "  test   Run all tests (pass -v for verbose)"
    echo "  build  Compile the Go binary"
    echo "  run    Execute the app (pass -add, -list, etc.)"
    exit 1
esac
