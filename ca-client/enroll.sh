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

CA_CLIENT="./bin/fabric-ca-client"
LOG="enroll.log"

> $LOG

# Override above defaults in ../config.env and optional local.env
source ../config.env
[ -r local.env ] && chmod og-rx local.env && source local.env

admin_server_uri=https://${CA_USER}:${CA_PASS}@ca-server.${CA_DOMAIN}:${CA_PORT}
orderer=${ORDERER_ORG}
peer=${PEER_ORG}

for i in orderer peer; do
    # ORDERER_ORG or PEER_ORG
    org="${!i}"

    admin_home="crypto-config/admin/${i}"

    export FABRIC_CA_CLIENT_TLS_CERTFILES="${PWD}/crypto-config/tlsca.${CA_DOMAIN}.pem"
    export FABRIC_CA_CLIENT_CANAME="ca-${i}-org"

    if [ ! -r ${admin_home}/msp/signcerts/cert.pem ]; then
	echo "Enrolling ${i} admin"
	${CA_CLIENT} enroll -u "${admin_server_uri}" -H "${admin_home}" \
	    2>>${LOG} || (tail -2 ${LOG}; exit 1)
    else
	echo "${i} admin already enrolled"
    fi

    name="Admin@${org}"
    home="crypto-config/${i}Organizations/${org}/users/$name"

    if [ ! -r $home/msp/signcerts/${name}-cert.pem ]; then
	# No ca-server auth needed, admin is enrolled
	server_uri="https://ca-server.${CA_DOMAIN}:${CA_PORT}"
	# Remove existing identity if there
	${CA_CLIENT} identity remove "${name}" -u "${server_uri}" -H "${admin_home}" >/dev/null 2>&1 || true

	# Set up id type/attrs
	case ${i} in
	  orderer)
	    id_args="--id.type=admin --id.attrs hf.Registrar.Roles=client,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,admin=true:ecert,abac.init=true:ecert"
	  ;;
	  *)
	    id_args="--id.type=user"
	  ;;
	esac

	echo "Registering ${name} with ${id_args}"
	# Secret appears on stdout of client register as "Password: <secret>"
	secret=$(
	    ${CA_CLIENT} register -u "${server_uri}" -H "${admin_home}" \
		--id.name "${name}" ${id_args} \
		2>>${LOG} | cut -f 2 -d ' '
	    ) || (tail -2 ${LOG}; exit 1)

	# Enroll name using secret from register
	csr_names="O=${org},C=${CSR_C},ST=${CSR_ST},L=${CSR_L}"
	server_uri="https://${name}:${secret}@ca-server.${CA_DOMAIN}:${CA_PORT}"

	echo "Enrolling ${name}"
	${CA_CLIENT} enroll -u "${server_uri}" -H "${home}" --csr.names "${csr_names}" \
	    2>>${LOG} || (tail -2 ${LOG}; exit 1)

	mv $home/msp/signcerts/cert.pem $home/msp/signcerts/${name}-cert.pem
    else
	echo "${name} already enrolled"
    fi

    orgmsp="crypto-config/${i}Organizations/${org}/msp"

    echo "Populating ${orgmsp} for channel genesis block"
    #mkdir -p $orgmsp/tlscacerts/
    #cp -f crypto-config/tlsca.${CA_DOMAIN}.pem $orgmsp/tlscacerts/
    mkdir -p $orgmsp/cacerts/ $orgmsp/admincerts/
    cp -f $home/msp/cacerts/*.pem $orgmsp/cacerts/
    cp -f $home/msp/signcerts/${name}-cert.pem $orgmsp/admincerts/
done
