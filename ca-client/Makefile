HLF_VERSION:=1.4.1

UNAME:=$(shell uname -s)
ARCH:=$(shell arch)
ARCHURL_Darwin-i386:=darwin-amd64
ARCHURL_Linux-x86_64:=linux-amd64
ARCHURL_CYGWIN_NT-10.0-x86_64:=windows-amd64
HLF_ARCH:=$(ARCHURL_$(UNAME)-$(ARCH))
ifndef HLF_ARCH
$(error Do not know how to handle $(UNAME)-$(ARCH))
endif
ifeq ($(UNAME), Darwin)
BASE64D=base64 -D
else
BASE64D=base64 -d
endif

include ../config.env

HLF_BINARY_URI=https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric-ca

ORDERER_TLSCA_PEM:=crypto-config/ordererOrganizations/$(ORDERER_ORG)/tlsca/tlsca.$(ORDERER_DOMAIN)-cert.pem
PEER_TLSCA_PEM:=crypto-config/peerOrganizations/$(PEER_ORG)/tlsca/tlsca.$(PEER_DOMAIN)-cert.pem

CA_SERVER_URI=https://ca-server.$(CA_DOMAIN):7054

.PHONY: all clean channel anchor-peers

all: $(ORDERER_TLSCA_PEM) $(PEER_TLSCA_PEM) crypto-config/.msp.stamp

crypto-config/bin/fabric-ca-client:
	@mkdir -p $(dir $@)
	@curl $(HLF_BINARY_URI)/hyperledger-fabric-ca/$(HLF_ARCH)-$(HLF_VERSION)/hyperledger-fabric-ca-$(HLF_ARCH)-$(HLF_VERSION).tar.gz | tar xz -C crypto-config

crypto-config/tlsca.$(CA_DOMAIN).pem:
	@mkdir -p $(dir $@)
	@curl -sk $(CA_SERVER_URI)/api/v1/cainfo | jq -r ".result.CAChain" | $(BASE64D) > $@ || rm -f $@
	@[ -s $@ ] || (rm -f $@; echo "Failed to get pem from $(CA_SERVER_URI)/api/v1/cainfo"; false)

%.pem: crypto-config/tlsca.$(CA_DOMAIN).pem
	@mkdir -p $(dir $@)
	@cp $< $@

crypto-config/.msp.stamp: crypto-config/bin/fabric-ca-client
	@./enroll.sh || (tail enroll.log; false)
	@touch $@

clean:
	rm -rf crypto-config artifacts bin

../config.env:
	$(MAKE) -C .. config.env

channel:

anchor-peers:
