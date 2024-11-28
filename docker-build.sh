#!/bin/bash

# 构建方式: 本地构建 or 多平台构建打包上传镜像
BUILD_MODE="local"  # local / remote

# DockerHub 用户名、镜像名称和版本号
USERNAME="ecnunic"
IMAGE_NAME="open-oauth2playground"
VERSION="v0.2.0"

# 支持的架构列表
PLATFORMS="linux/amd64,linux/arm64"

# 完整的镜像标签
FULL_TAG="${USERNAME}/${IMAGE_NAME}:${VERSION}"

echo "Building ${FULL_TAG} for platforms ${PLATFORMS}..."

if [ "${BUILD_MODE}" == "remote" ]; then
  # 推送到远程镜像仓库
  docker buildx build \
    --platform "${PLATFORMS}" \
    -t "${FULL_TAG}" \
    --push .
elif [ "${BUILD_MODE}" == "local" ]; then
  # 本地构建
  docker build --no-cache --load -t "${FULL_TAG}" .
else
  # Unknown $BUILD_MODE
  echo "BUILD_MODE must be \`local\` or \`remote\`, but got \`${BUILD_MODE}\`"
fi