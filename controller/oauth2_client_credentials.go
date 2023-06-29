package controller

import (
	"github.com/gin-gonic/gin"
)

type ReqCredentialData struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func clientCredentials(c *gin.Context) {

}
