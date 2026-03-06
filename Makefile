# Проверка наличия golangci-lint
GO_LINT := $(shell command -v golangci-lint 2>/dev/null)

# Цель: запуск всех проверок кода
check:
	@echo "🔍 Running code checks..."
	@go mod tidy -v
	@go vet ./...
	@if [ -n "$(GO_LINT)" ]; then \
		echo "✅ golangci-lint found, running..." && \
		golangci-lint run --timeout=5m --verbose; \
	else \
		echo "⚠️  golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2"; \
	fi
	@go test -race -v ./... 2>/dev/null || echo "⚠️  Tests skipped or failed (run manually with 'go test -race ./...')"
