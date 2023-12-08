package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ECNU/Open-OAuth2Playground/g"
	"github.com/ECNU/Open-OAuth2Playground/models"
	"github.com/gin-gonic/gin"
)

type ReqPkceData struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	CodeVerifier string `json:"code_verifier"`
	Scope        string `json:"scope"`
	RedirectURI  string `json:"redirect_uri"`
}

func pkce(c *gin.Context) {
	request := ReqPkceData{}
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}

	method := "POST"
	apiAddr := g.Config().Endpoints.Token
	grant_type := "authorization_code"
	body := fmt.Sprintf("code=%s&redirect_uri=%s&client_id=%s&scope=%s&grant_type=%s&code_verifier=%s",
		request.Code, request.RedirectURI, request.ClientID, request.Scope, grant_type, request.CodeVerifier)

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
