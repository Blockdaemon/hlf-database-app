include config.env

# include overrides if the file exists
-include local.env

ifndef GOPATH
    GOPATH:=$(HOME)/go
    export GOPATH
endif

MKFILES:=Makefile config.env $(wildcard local.env)	# only care about local.env if it is there
CHANFILE:=$(SERVICE_NETWORK)/artifacts/$(CHANNEL).channel.tx

.PHONY: all fmt
all: hlf-database-app config.yaml $(CHANFILE)

$(CHANFILE):
	make -C $(SERVICE_NETWORK) channel

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
	rm -rf __pycache__
	rm -f config.yaml
	rm -f hlf-database-app

.PHONY: clean-cc
clean-cc:
	docker ps -a | grep "chaincode -peer" | cut -f 1 -d " " | xargs -r docker rm
	docker image ls | grep "chaincode -peer" | cut -f 1 -d " " | xargs -r docker rmi

.PHONY: FORCE
