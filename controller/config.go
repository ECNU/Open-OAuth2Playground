package controller

import (
	"net/http"

	"github.com/ECNU/Open-OAuth2Playground/g"
	"github.com/gin-gonic/gin"
)

type ConfigData struct {
	TrustDomain           []string `json:"trust_domain"`
	AuthorizationEndpoint string   `json:"authorization_endpoint"`
	TokenEndpoint         string   `json:"token_endpoint"`
	UserinfoEndpoint      string   `json:"userinfo_endpoint"`
	DefaultScope          string   `json:"default_scope"`
}

func getConfig(c *gin.Context) {
	res := ConfigData{
		TrustDomain:           g.Config().TrustDomain,
		AuthorizationEndpoint: g.Config().Endpoints.Authorization,
		TokenEndpoint:         g.Config().Endpoints.Token,
		UserinfoEndpoint:      g.Config().Endpoints.Userinfo,
		DefaultScope:          g.Config().DefaultScope,
	}
	c.JSON(http.StatusOK, handleSuccess(res))
}
