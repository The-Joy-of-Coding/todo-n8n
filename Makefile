test:
	@INSECURE_API=http://localhost SECURE_API=https://example.com go test ./...

run:
	@go run .

build:
	@go build .
