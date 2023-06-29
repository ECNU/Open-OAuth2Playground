package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ECNU/Open-OAuth2Playground/g"
	"github.com/ECNU/Open-OAuth2Playground/models"
	"github.com/gin-gonic/gin"
)

type ReqCredentialData struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func clientCredentials(c *gin.Context) {
	reqData := ReqCredentialData{}
	if err := c.Bind(&reqData); err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}

	method := "POST"

	apiAddr := g.Config().Endpoints.Token

	grant_type := "client_credentials"
	clientId := reqData.ClientId
	clientSecret := reqData.ClientSecret
	body := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s", grant_type, clientId, clientSecret)

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
