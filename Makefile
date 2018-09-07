include config.env

# include overrides if the file exists
-include local.env

ifndef GOPATH
    GOPATH:=$(HOME)/go
    export GOPATH
endif

MAKEFILES:=Makefile config.env $(wildcard local.env)	# only care about local.env if it is there

.PHONY: all
all: hlf-database-app config.yaml

hlf-database-app: FORCE
	go build

config.yaml: $(MAKEFILES)

# jinja2 rule
%.yaml: templates/%.yaml.in $(MAKEFILES)
	NETWORK=$(NETWORK) DOMAIN=$(DOMAIN) CHANNEL=$(CHANNEL) CRYPTO=$(CRYPTO) tools/jinja2-cli.py < $< > $@ || (rm -f $@; false)

.PHONY: clean
clean:
	rm -rf __pycache__
	rm -f config.yaml
	rm -f hlf-database-app

.PHONY: FORCE
