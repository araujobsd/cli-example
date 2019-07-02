GO ?= go
SRC := main.go

all: build commands

build:
	@for file in $(SRC); do \
		go build -o thecarousell;\
	done \

commands:
	$(MAKE) -C ./commands

.PHONY: all build commands
