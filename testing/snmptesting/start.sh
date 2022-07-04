#!/bin/bash
export DOCKERFILE = ..\serviceswatcher\startup\docker-compose.yml
docker-compose -f ${DOCKERFILE} build --no-cache
docker-compose -f ${DOCKERFILE} up