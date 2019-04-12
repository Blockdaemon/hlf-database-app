#!/bin/bash

set -e

# "sane" defaults
CA_USER=admin
CA_PASS=adminpw
CA_DOMAIN=local
CA_PORT=7054

ORDERER_ORG="OrdererOrg"
PEER_ORG="PeerOrg"

CSR_C="US"
CSR_ST="California"
CSR_L="Los Angeles"

CA_CLIENT="./crypto-config/bin/fabric-ca-client"
LOG="enroll.log"

# Override above defaults in ../config.env and optional local.env
source ../config.env
[ -r local.env ] && chmod og-rx local.env && source local.env

admin_server_uri=https://${CA_USER}:${CA_PASS}@ca-server.${CA_DOMAIN}:${CA_PORT}
admin_home="crypto-config/admin"

export FABRIC_CA_CLIENT_TLS_CERTFILES="${PWD}/crypto-config/tlsca.${CA_DOMAIN}.pem"
export FABRIC_CA_CLIENT_CANAME="ca-server"

echo Enrolling admin
${CA_CLIENT} enroll -d -u "${admin_server_uri}" -H "${admin_home}" \
    2>${LOG} || (tail ${LOG}; exit 1)

orderer=${ORDERER_ORG}
peer=${PEER_ORG}
for i in orderer peer; do
    # ORDERER_ORG or PEER_ORG
    org="${!i}"
    name="Admin@${org}"

    # No ca-server auth needed, admin is enrolled
    server_uri="https://ca-server.${CA_DOMAIN}:${CA_PORT}"
    # Remove existing identity if there
    ${CA_CLIENT} identity remove "${name}" -u "${server_uri}" -H "${admin_home}" >/dev/null 2>&1 || true

    echo "Registering ${name}"
    # Secret appears on stdout of client register as "Password: <secret>"
    secret=$(
        ${CA_CLIENT} register -d -u "${server_uri}" -H "${admin_home}" \
	    --id.name "${name}" --id.secret "${secret}" --id.type=user \
	    2>>${LOG} | cut -f 2 -d ' '
	) || (tail ${LOG}; exit 1)

    # Enroll name using secret from register
    csr_names="O=$org,C=${CSR_C},ST=${CSR_ST},L=${CSR_L}"
    server_uri="https://${name}:${secret}@ca-server.${CA_DOMAIN}:${CA_PORT}"
    home="crypto-config/${i}Organizations/$org/users/$name"

    echo "Enrolling ${name}"
    ${CA_CLIENT} enroll -d -u "${server_uri}" -H "${home}" --csr.names "${csr_names}" \
	2>>${LOG} || (tail ${LOG}; exit 1)

    mv $home/msp/signcerts/cert.pem $home/msp/signcerts/${name}-cert.pem
done
