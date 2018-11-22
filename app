#!/bin/bash

if [ -z "${GOPATH}" ]; then
    export GOPATH=${HOME}/go
fi

source config.env
DOMAIN=${DOMAIN} ORG=${ORG} CHANNEL=${CHANNEL} ARTIFACTS=${ARTIFACTS} exec ./hlf-database-app "$@"
