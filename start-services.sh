#!/bin/bash

# author: dbbDylan
# date: 2024.11.09
# note: depends on `jq`

set -e  # 监测到错误立即退出

# ========================
# 变量定义
# ========================

# 通过环境变量初始化配置文件功能开关
CFG_INIT_ENABLE=${CFG_INIT_ENABLE:-1}

# docker 容器中各（配置）文件以及目录的路径
PATH_ROOT=${PATH_ROOT:-"/oauth2playground"}
PLAYGROUND_PATH="${PATH_ROOT}/oauth2playground"
PLAYGROUND_CONFIG_FILE="${PATH_ROOT}/cfg.json"

# 可对外暴露的环境变量
PLAYGROUND_PORT=${PLAYGROUND_PORT:-"80"}  # oauth2playground 服务端口号
OAUTH_SERVER_PORT=${OAUTH_SERVER_PORT:-"8081"}  # oauth-server-lite 服务端口号
OAUTH_SERVER_HOST=${OAUTH_SERVER_HOST:-"localhost"}  # oauth-server-lite 服务地址/域名
OAUTH_SERVER_URL=${OAUTH_SERVER_URL:-"http://${OAUTH_SERVER_HOST}:${OAUTH_SERVER_PORT}"}  # oauth-server-lite 服务 URL

# ========================
# 函数定义
# ========================

# 配置 oauth2playground cfg.json
configure_oauth2_playground() {
  echo "Configuring OAuth2 Playground..."

  # 更新 .endpoints 中的指定字段
  jq --arg url "$OAUTH_SERVER_URL" '
    .endpoints.authorization = "\($url)/oauth2/authorize" |
    .endpoints.device_authorization = "\($url)/oauth2/device/authorize" |
    .endpoints.token = "\($url)/oauth2/token" |
    .endpoints.userinfo = "\($url)/oauth2/userinfo"
  ' "$PLAYGROUND_CONFIG_FILE" > "$PLAYGROUND_CONFIG_FILE.tmp" && mv "$PLAYGROUND_CONFIG_FILE.tmp" "$PLAYGROUND_CONFIG_FILE"

  # 更新 .http 字段
  jq --arg port "$PLAYGROUND_PORT" \ '
    .http.listen = "0.0.0.0:\($port)"
  ' "$PLAYGROUND_CONFIG_FILE" > "$PLAYGROUND_CONFIG_FILE.tmp" && mv "$PLAYGROUND_CONFIG_FILE.tmp" "$PLAYGROUND_CONFIG_FILE"

  # 仅在 trust_domain 中不存在时追加新值
  jq --arg new_domain "${OAUTH_SERVER_HOST}:${OAUTH_SERVER_PORT}" '
    if .trust_domain | index($new_domain) == null then
      .trust_domain += [$new_domain]
    else
      .
    end
  ' "$PLAYGROUND_CONFIG_FILE" > "$PLAYGROUND_CONFIG_FILE.tmp" && mv "$PLAYGROUND_CONFIG_FILE.tmp" "$PLAYGROUND_CONFIG_FILE"

  echo "OAuth2 Playground configured successfully!"
}

# 启动 OAuth2 Playground 服务
start_oauth2_playground() {
  echo "Starting OAuth2 Playground..."

  cd "${PATH_ROOT}" && ${PLAYGROUND_PATH} -c "${PLAYGROUND_CONFIG_FILE}" &
}

# ========================
# 主执行流程
# ========================
# 判断 CFG_INIT_ENABLE 是否为 1
if [ "$CFG_INIT_ENABLE" -eq 1 ]; then
    configure_oauth2_playground
fi
start_oauth2_playground

# 保持脚本运行
echo "All services started. Keeping script running..."
tail -f /dev/null