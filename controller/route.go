package controller

import (
	"net/http"

	"github.com/ECNU/Open-OAuth2Playground/g"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	r.Static(g.Config().Http.RouteBase+"css/", "front-standalone/dist/css/")
	r.Static(g.Config().Http.RouteBase+"js/", "front-standalone/dist/js/")
	r.StaticFile(g.Config().Http.RouteBase, "front-standalone/dist/index.html")
	r.StaticFile(g.Config().Http.RouteBase+"favicon.ico", "front-standalone/dist/favicon.ico")

	r.GET(g.Config().Http.RouteBase+"v1/config", getConfig)

	playground := r.Group(g.Config().Http.RouteBase + "v1")
	playground.Use(IPLimitCheck)
	playground.Use(NoCache())
	playground.POST("/oauth2/pkce", pkce)
	playground.POST("/oauth2/device_flow", deviceFlow)
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
