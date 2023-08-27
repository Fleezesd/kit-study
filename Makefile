# ==============================================================================
# Includes

# 确保 `include common.mk` 位于第一行，common.mk 中定义了一些变量，后面的子 makefile 有依赖
include scripts/make-rules/common.mk 
include scripts/make-rules/generate.mk
include scripts/make-rules/tools.mk

.PHONY: protoc
protoc: ## 编译 protobuf 文件.
	@$(MAKE) gen.protoc