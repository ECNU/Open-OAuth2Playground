package controller

import (
	"github.com/gin-gonic/gin"
)

type ReqApiData struct {
	Method          string            `json:"method"`
	ApiAddr         string            `json:"api_addr"`
	AccessToken     string            `json:"access_token"`
	AccessTokenType string            `json:"access_token_type"`
	Header          map[string]string `json:"header"`
	HttpBody        string            `json:"http_body"`
}

func api(c *gin.Context) {
}
