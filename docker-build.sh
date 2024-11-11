#!/bin/bash

# DockerHub 用户名、镜像名称和版本号
USERNAME="ecnunic"
IMAGE_NAME="open-oauth2playground"
VERSION="v0.2.0"

# 支持的架构列表
PLATFORMS="linux/amd64,linux/arm64"

# 完整的镜像标签
FULL_TAG="${USERNAME}/${IMAGE_NAME}:${VERSION}"

echo "Building ${FULL_TAG} for platforms ${PLATFORMS}..."

# 推送到远程镜像仓库
docker buildx build \
    --platform "${PLATFORMS}" \
    -t "${FULL_TAG}" \
    --push .