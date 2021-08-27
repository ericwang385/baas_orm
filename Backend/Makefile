GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all:
	test build

test :
	$(GOTEST) codegen_test.go

build:
	$(GOBUILD)