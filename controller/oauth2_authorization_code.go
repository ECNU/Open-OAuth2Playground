package controller

import (
	"github.com/gin-gonic/gin"
)

type ExchangeTokenByCodeRequest struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	RedirectURI  string `json:"redirect_uri"`
}

type RefreshTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

func exchangeTokenByCode(c *gin.Context) {
}

func refreshToken(c *gin.Context) {
}
