GO ?= go
SRC := create_listing.go delete_listing.go get_category.go \
	get_listing.go register.go get_top_category.go \
	update_listing.go

all: build

build:
	@for file in $(SRC); do \
		go build $$file ;\
	done \

.PHONY: all build
