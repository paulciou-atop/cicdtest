#!/bin/bash
if [ -z "$1" ] 
then 
    DOCKER_PREFIX="tmp$(date '+%d%m%H%M')"
else 
    DOCKER_PREFIX=$1
fi
echo "DOCKER_PREFIX=${DOCKER_PREFIX}"
cp cicd/docker/Dockerfile_build_linux .
cp cicd/docker/docker-compose-cicdtest.yml . 
DOCKER_PREFIX=${DOCKER_PREFIX} docker compose -f docker-compose-cicdtest.yml build 
DOCKER_PREFIX=${DOCKER_PREFIX} docker compose -f docker-compose-cicdtest.yml up --detach
rm Dockerfile_build_linux
rm docker-compose-cicdtest.yml  