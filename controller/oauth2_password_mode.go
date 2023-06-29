package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ECNU/Open-OAuth2Playground/g"
	"github.com/ECNU/Open-OAuth2Playground/models"
	"github.com/gin-gonic/gin"
)

type ReqPasswordData struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

func passwordMode(c *gin.Context) {
	reqData := ReqPasswordData{}
	if err := c.Bind(&reqData); err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}

	method := "POST"

	apiAddr := g.Config().Endpoints.Token

	grant_type := "password"
	clientId := reqData.ClientId
	clientSecret := reqData.ClientSecret
	username := reqData.Username
	password := reqData.Password
	body := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s", grant_type, clientId, clientSecret, username, password)

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
