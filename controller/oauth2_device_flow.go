package controller

import (
	"fmt"
	"github.com/ECNU/Open-OAuth2Playground/g"
	"github.com/ECNU/Open-OAuth2Playground/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type ReqDeviceData struct {
	ClientId     string `json:"client_id"`
	Code         string `json:"code"`
	ResponseType string `json:"response_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type ReqUserCodeData struct {
	InitialAddress string `json:"initialAddress"`
}

func getUserCode(c *gin.Context) {
	reqData := ReqUserCodeData{}
	if err := c.Bind(&reqData); err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}
	method := "POST"
	apiAddr := reqData.InitialAddress
	body := fmt.Sprintf("")
	header := make(map[string]string)

	res, err := models.HandleRequest(method, apiAddr, g.UserAgent, body, g.Config().Timeout, header)
	if err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, handleSuccess(res))
}

func deviceFlow(c *gin.Context) {
	reqData := ReqDeviceData{}
	if err := c.Bind(&reqData); err != nil {
		c.JSON(http.StatusOK, handleError(err.Error()))
		return
	}

	method := "POST"
	apiAddr := g.Config().Endpoints.Token

	Code := reqData.Code
	clientId := reqData.ClientId
	ResponseType := reqData.ResponseType
	ExpiresIn := reqData.ExpiresIn
	body := fmt.Sprintf("code=%s&client_id=%s&response_type=%s", Code, clientId, ResponseType)

	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	header["Content-Length"] = strconv.Itoa(len(body))

	timeoutChan := time.After(time.Duration(ExpiresIn) * time.Second) // 设置expires_in秒后超时
	dataChan := make(chan models.Resp)

	// 启动一个协程循环检测token是否可用
	go func() {
		for {
			res, err := models.HandleRequest(method, apiAddr, g.UserAgent, body, g.Config().Timeout, header)
			if err != nil {
				c.JSON(http.StatusOK, handleError(err.Error()))
				return
			}
			// 检查返回值中是否包含token
			if rawJsonMap, ok := res.RawJson.(map[string]interface{}); ok {
				if _, tokenAvailable := rawJsonMap["access_token"]; tokenAvailable {
					dataChan <- res
					return
				}
			} else {
				fmt.Println("res.RawJson不是map格式")
			}
			time.Sleep(1000 * time.Millisecond) // 每隔1秒轮询
		}
	}()

	// 已获取token或超时
	select {
	case res := <-dataChan:
		c.JSON(http.StatusOK, handleSuccess(res))
		//return res, nil
	case <-timeoutChan:
		c.JSON(http.StatusOK, handleError("user code expired"))
	}
}
