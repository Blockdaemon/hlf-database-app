#!/bin/bash

if [ -z "${GOPATH}" ]; then
    export GOPATH=${HOME}/go
fi

set -a
source config.env
exec ./hlf-database-app "$@"
