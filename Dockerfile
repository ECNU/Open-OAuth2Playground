# 阶段1：构建（使用轻量级的 Alpine 版本作为构建阶段）

## 2.1 前端构建
FROM node:16-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app/Open-OAuth2Playground

# 复制前端项目并构建
COPY ./front-standalone /app/Open-OAuth2Playground/front-standalone
RUN cd front-standalone && npm install && npm run build


FROM golang:1.20-alpine AS backend-builder

# 设置工作目录
WORKDIR /app/Open-OAuth2Playground

# 复制依赖文件并安装依赖
COPY go.mod .
COPY go.sum .
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

# 复制源代码并编译
COPY . .
RUN CGO_ENABLED=0 go build -o oauth2playground .


# 阶段2：运行
FROM alpine:latest

# 设置工作目录并复制二进制文件
WORKDIR /app

ENV PATH_ROOT=/app

# 安装必要的运行 / 调试工具
RUN apk update && \
    apk add --no-cache sudo bash lsof jq curl iproute2 net-tools procps ca-certificates git iputils

COPY --from=frontend-builder /app/Open-OAuth2Playground/front-standalone/dist /app/front-standalone/dist
COPY --from=backend-builder /app/Open-OAuth2Playground/oauth2playground .
COPY --from=backend-builder /app/Open-OAuth2Playground/cfg-docker.json cfg.json

# 复制启动脚本
COPY start-services.sh ./start-services.sh

# 修改文件权限
RUN chmod +x ./oauth2playground
RUN chmod +x ./start-services.sh

EXPOSE 80

ENTRYPOINT ["./start-services.sh"]