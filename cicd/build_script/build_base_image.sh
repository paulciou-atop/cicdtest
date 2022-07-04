#!/bin/bash
cp cicd/docker/Dockerfile_base_linux .
docker build -t nms_base_image -f Dockerfile_base_linux . 
rm Dockerfile_base_linux
