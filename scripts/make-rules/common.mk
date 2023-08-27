# ==============================================================================
# 定义全局 Makefile 变量方便后面引用

SHELL := /bin/bash

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../../ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR := $(ROOT_DIR)/_output

# 定义包名
ROOT_PACKAGE=github.com/marmotedu/miniblog

# Protobuf 文件存放路径
APIROOT=$(ROOT_DIR)/pkg/proto

ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif
