# ==============================================================================
# 用来进行代码生成的 Makefile
#

.PHONY: gen.protoc
gen.protoc: tools.verify.protoc-gen-go ## 编译 protobuf 文件.
	@echo "===========> Generate protobuf files"
	@protoc                                            \
		--proto_path=$(APIROOT)                          \
		--proto_path=$(ROOT_DIR)/third_party             \
		--go_out=paths=source_relative:$(APIROOT)        \
		--go-grpc_out=paths=source_relative:$(APIROOT)   \
		$(shell find $(APIROOT) -name *.proto)