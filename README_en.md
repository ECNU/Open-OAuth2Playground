# Open-OAuth2Playground
Open-OAuth2Playground is an open source version of OAuth2.0 Playground that mimics Google’s OAuth 2.0 Playground

It supports local out-of-the-box use, suitable for OAuth2.0 learners to test and learn the OAuth2.0 protocol.

When used for server-side deployment, it is also very suitable for synchronizing with the OAuth2 Server documentation, making it easy for third-party callers to quickly develop and debug.

[English](./README_en.md) | [中文](./README.md)

![](https://github.com/ECNU/Open-OAuth2Playground/blob/main/demo.png?raw=true)

- [Open-OAuth2Playground](#open-oauth2playground)
	- [Installation and Running](#installation-and-running)
		- [Binary Direct Running](#binary-direct-running)
			- [Linux](#linux)
			- [Windows](#windows)
		- [Systemctl Hosting](#systemctl-hosting)
		- [Compilation and Packaging](#compilation-and-packaging)
			- [Backend Compilation](#backend-compilation)
			- [Frontend Compilation](#frontend-compilation)
			- [Unified Packaging](#unified-packaging)
		- [Configuration](#configuration)
			- [Backend Configuration](#backend-configuration)
				- [Backend Configuration Description](#backend-configuration-description)
			- [Frontend Configuration](#frontend-configuration)
				- [Frontend Configuration Description](#frontend-configuration-description)
				- [Frontend Deployment](#frontend-deployment)
				- [Customize Frontend Menu](#customize-frontend-menu)
		- [API](#api)
		- [Acknowledgements](#acknowledgements)

## Installation and Running

### Binary Direct Running
#### Linux
Download the latest [release] package from [release](https://github.com/ECNU/Open-OAuth2Playground/releases), unzip it and run it directly.

```
mkdir Open-OAuth2Playground
cd Open-OAuth2Playground/
wget https://github.com/ECNU/Open-OAuth2Playground/releases/download/v0.1.1/Open-OAuth2Playground-linux-0.1.1.tar.gz
tar -zxvf Open-OAuth2Playground-0.1.1-linux-amd64.tar.gz
./control start
```
Visit port 80 of your server to use it.

#### Windows
If you only need to run tests on Windows, you can download `Open-OAuth2Playground-windows-0.1.1.zip` from [release], unzip it and run `Open-OAuth2Playground.exe`.

### Systemctl Hosting
Assuming deployment in the `/opt/Open-OAuth2Playground` directory, if deployed in other directories, modify the `WorkingDirectory` and `ExecStart` fields in `playground.service`.
```
cp playground.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable playground
systemctl start playground
```

### Compilation and Packaging
Get the project source code
```
git clone https://github.com/ECNU/Open-OAuth2Playground.git
```
#### Backend Compilation
```
cd Open-OAuth2Playground/
go mod tidy
go build
```
#### Frontend Compilation
```
cd front-standalone/
pnpm install
pnpm build
```
#### Unified Packaging
```
cd ..
chmod +x control
./control pack
```

### Configuration
#### Backend Configuration
Refer to `cfg.json.example`, create `cfg.jon` configuration file, modify configuration as needed.
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
##### Backend Configuration Description
| Configuration Item | Type | Description |
| --- | --- | --- |
| logger.dir | string | Log folder |
| logger.level | string | Log level |
| logger.keepHours | int | Log retention time |
| endpoints.authorization | string | OAuth2.0 authorization address |
| endpoints.token | string | OAuth2.0 get token address |
| endpoints.userinfo | string | OAuth2.0 get user information address |
| iplimit.enable | bool | Whether to enable IP restriction |
| iplimit.trust_ip | []string | List of trusted IP addresses |
| http.route_base | string | Route prefix, note to match with frontend |
| http.trust_proxy | []string | List of trusted proxy IP addresses |
| http.cors | []string | List of domain names allowed for frontend cross-domain access |
| http.listen | string | Listening address |
| trust_domain | []string | List of trusted domain names when backend forwards API calls |
| default_scope | string | Default scope |
| timeout | int | Timeout time |

#### Frontend Configuration
Modify `.env.production`
```ini
# Router path
VUE_APP_ROUTER_BASE=/
# Api Config
VUE_APP_API_PROTO=http
VUE_APP_API_HOST=localhost
VUE_APP_API_PORT=
VUE_APP_API_VERSION=v1
```
##### Frontend Configuration Description
| Configuration Item | Type| Description |
| --- | --- | --- |
| VUE_APP_ROUTER_BASE | string | Route prefix, note to match with backend |
| VUE_APP_API_PROTO | string | Required when frontend is deployed independently, the proto of the backend server |
| VUE_APP_API_HOST | string | Required when frontend is deployed independently, the domain name of the backend |
| VUE_APP_API_PORT | string | Required when frontend is deployed independently, the port of the backend. If it is the default port, it can be ignored (such as https'443 or http'80) |
| VUE_APP_API_VERSION | string | API version, currently fixed to v1 |

##### Frontend Deployment
The frontend part of the project can be deployed independently or published by the backend.

By default, it is published by the backend, in which case the frontend's `VUE_APP_API_HOST`,`VUE_APP_API_PROTO`, `VUE_APP_API_PORT` and other configuration items can be ignored. At this time, the compiled and packaged frontend code should be deployed under the front-standalone/dist directory relative to the backend binary file.

If the frontend is deployed independently, you need to configure `VUE_APP_API_HOST`,`VUE_APP_API_PROTO`, `VUE_APP_API_PORT` and other configuration items at compile time, and make sure that the frontend domain name is in the backend's cross-domain list.

##### Customize Frontend Menu

The menu part of the project corresponds to the `front-standalone/src/views/Layourt.vue` file, you can modify the content of `el-menu-item` as needed, and then compile and package it.

### API
todo


### Acknowledgements
This project was inspired by Google's [OAuth 2.0 Playground](https://developers.google.com/oauthplayground/)

Thanks to Google for providing excellent tools.