#!/bin/bash
set -e

DOCKER_IMAGE_NAME="rdbell/gptswe"
DOCKER_TAG=$(date +%s)

# Build and push Docker image
docker buildx build --platform linux/amd64,linux/arm64 -t $DOCKER_IMAGE_NAME:$DOCKER_TAG --push .
docker buildx build --platform linux/amd64,linux/arm64 -t $DOCKER_IMAGE_NAME:latest --push .
