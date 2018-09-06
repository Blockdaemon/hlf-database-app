#!/bin/bash
sed -i "
s/heroes-service-network/{{env['NETWORK'] or \"some-network\"}}/
s/chainhero.io/{{env['DOMAIN'] or \"localhost\"}}/
s#\${GOPATH}/src/github.com/chainHero/heroes-service/fixtures/crypto-config#{{env['CRYPTO']}}#
" $1
