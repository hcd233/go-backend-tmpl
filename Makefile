# Makefile for aris-api-tmpl

APP_NAME := aris-api-tmpl
MAIN     := main.go
OUTPUT   := $(APP_NAME)

GOMAXPROCS ?= $(shell nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)

# 显式测试包列表：cmd + internal + test（不含 web 等非 Go 目录）
GO_TEST_PACKAGES := ./cmd/... ./internal/... ./test/...

# golangci-lint 版本（CI 与本地对齐）
GOLANGCI_LINT_VERSION := v2.12.2

LDFLAGS     := -s -w
BUILD_FLAGS := -trimpath -p $(GOMAXPROCS)

.PHONY: build build-upx build-dev build-debug warm-cache clean clean-all fmt vet test test-cover lint lint-static help

## build: 生产构建（strip 符号）
build:
	CGO_ENABLED=0 go build \
		$(BUILD_FLAGS) \
		-ldflags="$(LDFLAGS)" \
		-o $(OUTPUT) $(MAIN)
	@echo "Built $(OUTPUT) ($$(du -h $(OUTPUT) | cut -f1))"

## build-upx: 极致压缩构建（strip + UPX，体积最小，需安装 upx）
build-upx: build
	upx --best --lzma $(OUTPUT)
	@echo "Compressed $(OUTPUT) ($$(du -h $(OUTPUT) | cut -f1))"

## build-dev: 开发构建（保留调试信息，最快编译速度）
build-dev:
	go build -p $(GOMAXPROCS) \
		-o $(OUTPUT) $(MAIN)
	@echo "Built $(OUTPUT) ($$(du -h $(OUTPUT) | cut -f1))"

## build-debug: 带完整调试信息的构建（用于 dlv 调试）
build-debug:
	go build -p $(GOMAXPROCS) \
		-gcflags="all=-N -l" \
		-o $(OUTPUT) $(MAIN)
	@echo "Built $(OUTPUT) ($$(du -h $(OUTPUT) | cut -f1))"

## warm-cache: 预热编译缓存（CI 首次运行后可加速后续编译）
warm-cache:
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" -o /dev/null $(MAIN)
	@echo "Build cache warmed"

## clean: 清理构建产物
clean:
	rm -f $(OUTPUT)

## clean-all: 清理构建产物、测试覆盖率和编译缓存
clean-all: clean
	rm -f coverage.out coverage.html
	go clean -cache

## fmt: 格式化 Go 代码
fmt:
	go fmt ./...

## vet: 运行 go vet
vet:
	go vet ./...

## test: 运行全量测试（cmd + internal + test）
test:
	go test -count=1 $(GO_TEST_PACKAGES)

## test-cover: 带覆盖率的测试
test-cover:
	go test -count=1 -coverprofile=coverage.out $(GO_TEST_PACKAGES)
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: 运行全部静态检查
lint: lint-static

## lint-static: 运行 Go 静态分析（go vet + 可选 golangci-lint）
lint-static:
	@go run $(MAIN) lint static ./...

## help: 显示帮助
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | column -t -s ':'
