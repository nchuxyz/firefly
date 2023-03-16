.PHONY: build

GO       := GO111MODULE=on go
GOBUILD  := CGO_ENABLED=0 $(GO) build $(BUILD_FLAG)

default: build buildsucc

buildsucc:
	@echo Build FireFly Server successfully!

build:
	$(GOBUILD) -o bin/ffsvr cmd/server/main.go