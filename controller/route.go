package controller

import (
	"net/http"

	"github.com/ECNU/Open-OAuth2Playground/g"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	r.GET("/playground/v1/config", getConfig)

	playground := r.Group("/playground/v1")
	playground.Use(IPLimitCheck)
	playground.Use(NoCache())
	playground.POST("/oauth2/client_credentials", clientCredentials)
	playground.POST("/oauth2/password", passwordMode)
	playground.POST("/oauth2/authorization_code", exchangeTokenByCode)
	playground.POST("/oauth2/refresh_token", refreshToken)
	playground.POST("/api", api)
}

func InitGin(listen string) (httpServer *http.Server) {
	if g.Config().Logger.Level == "DEBUG" {
		gin.SetMode((gin.DebugMode))
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	if g.Config().Logger.Level == "DEBUG" {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())
	r.Use(CORS())

	Routes(r)
	httpServer = &http.Server{
		Addr:    g.Config().Http.Listen,
		Handler: r,
	}
	return
}
