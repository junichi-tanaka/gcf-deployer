
TARGET := gcf-deployer
BINDIR=./bin

VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS := -X 'main.version=$(VERSION)'

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(TARGET)

