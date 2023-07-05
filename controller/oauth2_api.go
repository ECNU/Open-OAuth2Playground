package controller

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ECNU/Open-OAuth2Playground/g"
	"github.com/ECNU/Open-OAuth2Playground/models"
	"github.com/gin-gonic/gin"
)

type ReqApiData struct {
	Method          string            `json:"method"`
	ApiAddr         string            `json:"api_addr"`
	AccessToken     string            `json:"access_token"`
	AccessTokenType string            `json:"access_token_type"`
	Header          map[string]string `json:"header"`
	HttpBody        string            `json:"http_body"`
}

func base64Decode(input string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	decodedString := string(decodedBytes)
	return decodedString, nil
}

func api(c *gin.Context) {
	reqData := ReqApiData{}
	if err := c.Bind(&reqData); err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}
	apiAddrUrl, err := url.Parse(reqData.ApiAddr)
	if err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}
	if !models.InSliceStr(apiAddrUrl.Host, g.Config().TrustDomain) {
		c.JSON(http.StatusOK, handleError("api addr not in trust domain"))
		return
	}

	method := reqData.Method
	apiAddr := reqData.ApiAddr
	header := reqData.Header
	body, err := base64Decode(reqData.HttpBody)
	if err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}
	header["Content-Length"] = strconv.Itoa(len(body))

	if reqData.AccessTokenType == "bearer" {
		header["Authorization"] = "Bearer " + reqData.AccessToken
	} else if reqData.AccessTokenType == "query" {
		apiAddr += "?access_token=" + reqData.AccessToken
	} else {
		c.JSON(http.StatusOK, handleError("invalid access token"))
		return
	}

	res, err := models.HandleRequest(method, apiAddr, g.UserAgent, body, g.Config().Timeout, header)
	if err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, handleSuccess(res))
}
