# 设置 Go 编译器
GO := go

# 设置项目名称和主程序路径
BINARY_NAME := nginxcert
MAIN_PATH := ./cmd/main.go

# 设置构建标志
BUILD_FLAGS := -v

# 默认目标
.PHONY: all
all: build

# 构建项目
.PHONY: build
build:
	$(GO) build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

# 运行测试
.PHONY: test
test:
	$(GO) test -v ./...

# 清理构建产物
.PHONY: clean
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

# 运行程序
.PHONY: run
run: build
	./$(BINARY_NAME)

# 运行程序（带域名过滤器）
.PHONY: run-filtered
run-filtered:
	@echo "请输入域名过滤器 (逗号分隔，例如 domain1.com,domain2.com):"
	@read DOMAIN_FILTER; \
	./$(BINARY_NAME) -config-path ./example/conf.d -author your@email.com -ssl-path ./example/ssl -domain-filter $$DOMAIN_FILTER

# 发布新版本
.PHONY: release
release:
	@echo "请输入版本号 (例如 v1.0.0):"
	@read VERSION; \
	./scripts/release.sh $$VERSION

# 帮助信息
.PHONY: help
help:
	@echo "可用的 make 命令:"
	@echo "  make build         - 构建项目"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理构建产物"
	@echo "  make run           - 构建并运行项目"
	@echo "  make run-filtered  - 构建并运行项目（带域名过滤器）"
	@echo "  make release       - 发布新版本"
	@echo "  make help          - 显示此帮助信息"