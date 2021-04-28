go_lint:
	docker run --rm -v ${PWD}:/app -w /app/ golangci/golangci-lint:v1.36-alpine golangci-lint run -v --timeout=5m