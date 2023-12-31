BUILDER=./.builder
RULES=go
include $(BUILDER)/rules.mk
$(BUILDER)/rules.mk:
	-go run github.com/endobit/builder@latest init

# sqlc

generate::
	sqlc generate

lint::
	sqlc compile

# protobuf

generate::
	buf generate
	go generate ./...

lint::
	cd proto && buf lint

nuke::
	rm -rf gen

format::
	buf format -w

# code

build::
	$(GO_BUILD) ./cmd/stackd
	$(GO_BUILD) ./cmd/stack

clean::
	rm -f stackd stack
