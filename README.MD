# Open-OAuth2Playground
Open-OAuth2Playground: 简单而好用的开源 OAuth2-Playground 

![](https://github.com/ECNU/Open-OAuth2Playground/blob/main/demo.jpg?raw=true)

- [Open-OAuth2Playground](#open-oauth2playground)
  - [安装运行](#安装运行)
    - [二进制直接运行](#二进制直接运行)
    - [systemctl 托管](#systemctl-托管)
    - [编译打包](#编译打包)
      - [后端编译](#后端编译)
      - [前端编译](#前端编译)
      - [统一打包](#统一打包)
    - [配置](#配置)
      - [后端配置](#后端配置)
      - [前端配置](#前端配置)
    - [API](#api)


## 安装运行

### 二进制直接运行
在 [release](https://github.com/ECNU/Open-OAuth2Playground/releases) 中下载最新的 [release] 包，解压后直接运行即可。

```
tar -zxvf Open-OAuth2Playground-0.1.0-linux-amd64.tar.gz
cd Open-OAuth2Playground/
./control start
```
访问你服务器的 80 端口即可使用。


### systemctl 托管
假定部署在 `/opt/Open-OAuth2Playground` 目录下，如果部署在其他目录修改 `playground.service` 中的 `WorkingDirectory` 和 `ExecStart` 两个字段即可。
```
cp playground.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable playground
systemctl start playground
```

### 编译打包
获取项目源码
```
git clone https://github.com/ECNU/Open-OAuth2Playground.git
```
#### 后端编译
```
cd Open-OAuth2Playground/
go mod tidy
go build
```
#### 前端编译
```
cd front-standalone/
pnpm install
pnpm build
```
#### 统一打包
```
cd ..
chmod +x control
./control pack
```

### 配置
#### 后端配置
修改 `cfg.json`

#### 前端配置
修改 `.env.production`

注意前端的 `route_base` 需要和后端相匹配。
### API
todo