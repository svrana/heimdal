SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

DATE 		:= $(shell date +"%a %b %d %T %Y")
UNAME_S 	:= $(shell uname -s | tr A-Z a-z)

build: ## Build heimdal
	go build -o heimdal cmd/heimdal/main.go

install:
	cp heimdal ${BIN_DIR}

.PHONY: help
help:  ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[/.a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

