package controller

import (
	"github.com/gin-gonic/gin"
)

type ReqPasswordData struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

func passwordMode(c *gin.Context) {

}
