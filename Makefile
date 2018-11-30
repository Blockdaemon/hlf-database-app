include config.env

# include overrides if the file exists
-include local.env

ifndef GOPATH
    GOPATH:=$(HOME)/go
    export GOPATH
endif

MKFILES:=Makefile config.env $(wildcard local.env)	# only care about local.env if it is there
CHANFILES:=$(SERVICE_NETWORK)/artifacts/$(CHANNEL).channel.tx $(SERVICE_NETWORK)/artifacts/$(CHANNEL).anchor-peers.tx
ifeq ($(shell uname -s),Darwin)
XARGS:=xargs
else
XARGS:=xargs -r
endif
.PHONY: all fmt
all: hlf-database-app config.yaml $(CHANFILES)

$(CHANFILES):
	make -C $(SERVICE_NETWORK) channel anchor-peers

fmt:
	gofmt -w $(wildcard *.go */*.go)

hlf-database-app: FORCE
	go build

config.yaml: $(MKFILES)

# jinja2 rule
%.yaml: templates/%.yaml.in $(MKFILES)
	NETWORK=$(NETWORK) DOMAIN=$(DOMAIN) CHANNEL=$(CHANNEL) CRYPTO=$(CRYPTO) tools/jinja2-cli.py < $< > $@ || (rm -f $@; false)

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
