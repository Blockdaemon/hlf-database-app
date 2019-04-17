DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# jinja2 rule
%.yaml: templates/%.yaml.in $(MKFILES)
	eval $$(sed -e 's/#.*$$//' $(DIR)/config.env) $(DIR)/tools/jinja2-cli.py < $< > $@ || (rm -f $@; false)
