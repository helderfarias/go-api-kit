PKG_LIST_ALL_TESTS	:= $(shell go list ./... | grep -v /vendor | grep -v /gopdf)

all: help

test:
	@go test -count=1 -cover $(PKG_LIST_ALL_TESTS)

help:
	@echo 'Usage: '
	@echo 'make test'

.PHONY: all test help
