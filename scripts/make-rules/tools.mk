# ==============================================================================
# 工具相关的 Makefile
#

TOOLS ?= protoc-gen-go 

.PHONY: tools.verify
tools.verify: $(addprefix tools.verify., $(TOOLS))

.PHONY: tools.install
tools.install: $(addprefix tools.install., $(TOOLS))

.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: install.protoc-gen-go
install.protoc-gen-go:
	@$(GO) install github.com/golang/protobuf/protoc-gen-go@latest