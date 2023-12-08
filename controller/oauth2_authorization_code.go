package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ECNU/Open-OAuth2Playground/g"
	"github.com/ECNU/Open-OAuth2Playground/models"
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
	ClientSecret string `json:"client_secret,omitempty"` // client_secret is optional for refresh_token
	RefreshToken string `json:"refresh_token"`
}

func exchangeTokenByCode(c *gin.Context) {
	request := ExchangeTokenByCodeRequest{}
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}

	method := "POST"
	apiAddr := g.Config().Endpoints.Token
	grant_type := "authorization_code"
	body := fmt.Sprintf("code=%s&redirect_uri=%s&client_id=%s&client_secret=%s&scope=%s&grant_type=%s", request.Code, request.RedirectURI, request.ClientID, request.ClientSecret, request.Scope, grant_type)

	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	header["Content-Length"] = strconv.Itoa(len(body))

	res, err := models.HandleRequest(method, apiAddr, g.UserAgent, body, g.Config().Timeout, header)
	if err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, handleSuccess(res))
}

func refreshToken(c *gin.Context) {
	request := RefreshTokenRequest{}
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}

	method := "POST"
	apiAddr := g.Config().Endpoints.Token
	grant_type := "refresh_token"
	body := fmt.Sprintf("grant_type=%s&client_id=%s&refresh_token=%s", grant_type, request.ClientID, request.RefreshToken)

	if request.ClientSecret != "" {
		body += fmt.Sprintf("&client_secret=%s", url.QueryEscape(request.ClientSecret))
	}

	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	header["Content-Length"] = strconv.Itoa(len(body))

	res, err := models.HandleRequest(method, apiAddr, g.UserAgent, body, g.Config().Timeout, header)
	if err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, handleSuccess(res))
}
