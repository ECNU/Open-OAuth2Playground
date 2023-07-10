# Open-OAuth2Playground
Open-OAuth2Playground: 开箱即用的开源版 OAuth2-Playground 

![](https://github.com/ECNU/Open-OAuth2Playground/blob/main/demo.png?raw=true)

- [Open-OAuth2Playground](#open-oauth2playground)
  - [安装运行](#安装运行)
    - [二进制直接运行](#二进制直接运行)
      - [Linux](#linux)
      - [Windows](#windows)
    - [systemctl 托管](#systemctl-托管)
    - [编译打包](#编译打包)
      - [后端编译](#后端编译)
      - [前端编译](#前端编译)
      - [统一打包](#统一打包)
    - [配置](#配置)
      - [后端配置](#后端配置)
        - [后端配置说明](#后端配置说明)
      - [前端配置](#前端配置)
        - [前端配置说明](#前端配置说明)
        - [前端部署](#前端部署)
    - [API](#api)


## 安装运行

### 二进制直接运行
#### Linux
在 [release](https://github.com/ECNU/Open-OAuth2Playground/releases) 中下载最新的 [release] 包，解压后直接运行即可。

```
mkdir Open-OAuth2Playground
cd Open-OAuth2Playground/
wget https://github.com/ECNU/Open-OAuth2Playground/releases/download/v0.1.1/Open-OAuth2Playground-linux-0.1.1.tar.gz
tar -zxvf Open-OAuth2Playground-0.1.1-linux-amd64.tar.gz
./control start
```
访问你服务器的 80 端口即可使用。

#### Windows
如果只需要在 Windows 上运行测试，可以直接下载 [release] 中的 `Open-OAuth2Playground-windows-0.1.1.zip`，解压后运行 `Open-OAuth2Playground.exe` 即可。

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
参考 `cfg.json.example`, 创建 `cfg.jon` 配置文件，修改配置即可
```json
{
	"logger": {
		"dir": "logs/",
		"level": "DEBUG",
		"keepHours": 24
	},
	"endpoints": {
		"authorization": "http://cas.example.org/cas/oauth2.0/authorize",
		"token": "http://cas.example.org/cas/oauth2.0/accessToken",
		"userinfo": "http://cas.example.org/cas/oauth2.0/profile"
	},
	"iplimit": {
		"enable": false,
		"trust_ip": ["127.0.0.1","::1"]
	},
	"http": {
		"route_base":"/",
		"trust_proxy": ["127.0.0.1", "::1"],
		"cors": ["http://127.0.0.1:8080","http://localhost:8080"],
		"listen": "0.0.0.0:80"
	},
	"trust_domain": ["cas.example.org", "localhost"],
	"default_scope": "Basic",
	"timeout": 10
}
```
##### 后端配置说明
| 配置项 | 类型 | 说明 |
| --- | --- | --- |
| logger.dir | string | 日志文件夹 |
| logger.level | string | 日志等级 |
| logger.keepHours | int | 日志保留时间 |
| endpoints.authorization | string | OAuth2.0 授权地址 |
| endpoints.token | string | OAuth2.0 获取 token 地址 |
| endpoints.userinfo | string | OAuth2.0 获取用户信息地址 |
| iplimit.enable | bool | 是否开启 IP 限制 |
| iplimit.trust_ip | []string | 可信任的 IP 列表 |
| http.route_base | string | 路由前缀，注意要和前端匹配 |
| http.trust_proxy | []string | 可信任的代理 IP 列表 |
| http.cors | []string | 允许前端跨域的域名列表 |
| http.listen | string | 监听地址 |
| trust_domain | []string | 后端转发API调用时，信任的域名列表 |
| default_scope | string | 默认的 scope |
| timeout | int | 超时时间 |


#### 前端配置
修改 `.env.production`
```ini
# Router path
VUE_APP_ROUTER_BASE=/
# Api Config
VUE_APP_API_PROTO=http
VUE_APP_API_HOST=localhost
VUE_APP_API_PORT=
VUE_APP_API_VERSION=v1
```
##### 前端配置说明
| 配置项 | 类型| 说明 |
| --- | --- | --- |
| VUE_APP_ROUTER_BASE | string | 路由前缀，注意要和后端匹配 |
| VUE_APP_API_PROTO | string | 前端独立部署时需要，后端服务器的 proto |
| VUE_APP_API_HOST | string | 前端独立部署时需要，后端的域名 |
| VUE_APP_API_PORT | string | 前端独立部署时需要，后端的端口。如果是默认端口可忽略（例如https的443或者http的80） |
| VUE_APP_API_VERSION | string | API 版本，目前固定为 v1 |

##### 前端部署
项目的前端部分可以独立部署单独发布，也可以由后端来发布。

默认由后端发布，此时前端的 `VUE_APP_API_HOST`,`VUE_APP_API_PROTO`, `VUE_APP_API_PORT` 等配置项可以忽略。此时编译后的前端代码应该部署相对后端二进制文件的 front-standalone/dist 目录下。

如果前端独立部署，则需要在编译时，配置好 `VUE_APP_API_HOST`,`VUE_APP_API_PROTO`, `VUE_APP_API_PORT` 等配置项，并确保前端的域名在后端的跨域列表内。
### API
todo