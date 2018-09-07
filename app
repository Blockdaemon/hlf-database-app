#!/bin/bash
source config.env
DOMAIN=${DOMAIN} CHANNEL=${CHANNEL} ARTIFACTS=${ARTIFACTS} exec ./hlf-database-app "$@"
