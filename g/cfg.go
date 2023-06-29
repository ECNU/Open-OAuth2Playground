package g

import (
	"encoding/json"

	"log"
	"sync"

	"github.com/toolkits/file"
)

/*
GlobalConfig 全局配置
*/
type GlobalConfig struct {
	Logger       LoggerSection  `json:"logger"`
	Endpoints    EndpointConfig `json:"endpoints"`
	IpLimit      IpLimitConfig  `json:"iplimit"`
	Http         HttpConfig     `json:"http"`
	TrustDomain  []string       `json:"trust_domain"`
	DefaultScope string         `json:"default_scope"`
	Timeout      int            `json:"timeout"`
}

type IpLimitConfig struct {
	Enable  bool
	TrustIP []string `json:"trust_ip"`
}

/*
EndpointConfig oauth endpoint 配置
*/
type EndpointConfig struct {
	Authorization string `json:"authorization"`
	Token         string `json:"token"`
	Userinfo      string `json:"userinfo"`
}

/*
HttpConfig Http 配置
*/
type HttpConfig struct {
	Listen     string   `json:"listen"`
	CORS       []string `json:"cors"`
	TrustProxy []string `json:"trust_proxy"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

/*
Config 安全的读取和修改配置
*/
func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

/*
ParseConfig 加载配置
*/
func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c
}
