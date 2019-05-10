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
# /usr/local/bin/base64 might support -d but if it does it won't like -D
BASE64D=/usr/bin/base64 -D
else
BASE64D=base64 -d
endif

include ../config.env

HLF_CA_BINARY_URI:=https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric-ca/hyperledger-fabric-ca
HLF_BINARY_URI:=https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric

CA_SERVER_URI:=https://ca-server.$(CA_DOMAIN):7054

CRYPTO_ORDERER_DIR:=crypto-config/ordererOrganizations/$(ORDERER_ORG)
CRYPTO_PEER_DIR:=crypto-config/peerOrganizations/$(PEER_ORG)

ORDERER_TLSCA_PEM:=$(CRYPTO_ORDERER_DIR)/tlsca/tlsca.$(ORDERER_DOMAIN)-cert.pem
PEER_TLSCA_PEM:=$(CRYPTO_PEER_DIR)/tlsca/tlsca.$(PEER_DOMAIN)-cert.pem

ARTIFACTS:=artifacts/$(CHANNEL).channel.tx artifacts/$(CHANNEL).anchor-peers.tx

CRYPTO_STAMP:=crypto-config/.msp.stamp

.PHONY: all bin-clean dist-clean clean

all: $(ORDERER_TLSCA_PEM) $(PEER_TLSCA_PEM) $(ARTIFACTS)

bin/fabric-ca-client:
	@mkdir -p $(dir $@)
	@curl $(HLF_CA_BINARY_URI)/$(HLF_ARCH)-$(HLF_VERSION)/hyperledger-fabric-ca-$(HLF_ARCH)-$(HLF_VERSION).tar.gz | tar xz

bin/configtxgen:
	@mkdir -p $(dir $@)
	@curl $(HLF_BINARY_URI)/$(HLF_ARCH)-$(HLF_VERSION)/hyperledger-fabric-$(HLF_ARCH)-$(HLF_VERSION).tar.gz | tar xz

crypto-config/tlsca.$(CA_DOMAIN).pem:
	@mkdir -p $(dir $@)
	@curl -sk $(CA_SERVER_URI)/api/v1/cainfo | jq -r ".result.CAChain" | $(BASE64D) > $@ || rm -f $@
	@[ -s $@ ] || (rm -f $@; echo "Failed to get pem from $(CA_SERVER_URI)/api/v1/cainfo"; false)

%.pem: crypto-config/tlsca.$(CA_DOMAIN).pem
	@mkdir -p $(dir $@)
	@cp $< $@

$(CRYPTO_STAMP): bin/fabric-ca-client
	@./enroll.sh || (tail enroll.log; false)
	@touch $@

artifacts/$(CHANNEL).channel.tx: bin/configtxgen configtx.yaml $(CRYPTO_STAMP)
	@mkdir -p artifacts
	bin/configtxgen -profile $(PROFILE)Channel -outputCreateChannelTx $@ -channelID $(CHANNEL)

artifacts/$(CHANNEL).anchor-peers.tx: bin/configtxgen configtx.yaml
	@mkdir -p artifacts
	bin/configtxgen -profile $(PROFILE)Channel -outputAnchorPeersUpdate $@ -channelID $(CHANNEL) -asOrg $(PEER_ORGNAME)

clean:
	rm -rf crypto-config artifacts

bin-clean:
	rm -rf bin

dist-clean: clean bin-clean

../config.env:
	$(MAKE) -C .. config.env

-include ../rules.mk
