package controller

import (
	"net/http"

	"github.com/ECNU/Open-OAuth2Playground/g"
	"github.com/ECNU/Open-OAuth2Playground/models"
	"github.com/gin-gonic/gin"
)

type GlobalResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func handleError(msg string) (res GlobalResponse) {
	res.Code = -1
	res.Msg = msg
	return
}

func handleSuccess(data interface{}) (res GlobalResponse) {
	res.Code = 0
	res.Msg = "ok"
	res.Data = data
	return
}

func IPLimitCheck(c *gin.Context) {
	if g.Config().IpLimit.Enable {
		if !models.IPCheck(c.ClientIP(), g.Config().IpLimit.TrustIP) {
			c.JSON(http.StatusOK, handleError("Your IP is not allowed"))
			c.Abort()
			return
		}
	}
	c.Next()
}

func CORS() gin.HandlerFunc {
	return func(context *gin.Context) {
		if models.InSliceStrFuzzy(context.Request.Header.Get("Origin"), g.Config().Http.CORS) {
			context.Writer.Header().Add("Access-Control-Allow-Origin", context.Request.Header.Get("Origin"))
			context.Writer.Header().Set("Access-Control-Max-Age", "86400")
			context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Cookie")
			context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		}

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}

func NoCache() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Add("Cache-Control", "no-store")
		context.Writer.Header().Add("Pragma", "no-cache")
		context.Next()
	}
}
