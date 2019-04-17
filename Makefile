include config.env

ifndef GOPATH
    GOPATH:=$(HOME)/go
    export GOPATH
endif

MKFILES:=Makefile config.env $(wildcard local.env)	# only care about local.env if it is there
CHANFILES:=$(ARTIFACTS)/$(CHANNEL).channel.tx $(ARTIFACTS)/artifacts/$(CHANNEL).anchor-peers.tx
ifeq ($(shell uname -s),Darwin)
XARGS:=xargs
else
XARGS:=xargs -r
endif
.PHONY: all fmt
all: hlf-database-app config.yaml $(CHANFILES)
config.env:
	@cp examples/config.env config.env

$(CHANFILES): FORCE
	make -C $(WORK_DIR) channel anchor-peers

fmt:
	gofmt -w $(wildcard *.go */*.go)

hlf-database-app: FORCE
	go build

config.yaml: $(MKFILES)

.PHONY: clean
clean:
	go clean
	rm -rf __pycache__
	rm -f config.yaml
	rm -f hlf-database-app

.PHONY: clean-cc
clean-cc:
	docker ps -a | grep "hlf-database-app" | cut -f 1 -d " " | $(XARGS) docker rm
	docker image ls | grep "hlf-database-app" | cut -f 1 -d " " | $(XARGS) docker rmi

.PHONY: FORCE
-include rules.mk
