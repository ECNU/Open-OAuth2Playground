package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_httpRequest(t *testing.T) {
	method := "GET"
	url := "https://ip.ecnu.edu.cn/myip"
	userAgent := "Open-OAuth2Playground"
	timeout := 10
	req, resp, err := httpRequest(method, url, userAgent, timeout, nil, nil)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, req, nil)
	assert.NotEqual(t, resp, nil)
}

func Test_handleRequest(t *testing.T) {
	method := "GET"
	url := "https://ip.ecnu.edu.cn/myip"
	userAgent := "Open-OAuth2Playground"
	timeout := 10
	res, err := HandleRequest(method, url, userAgent, "", timeout, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, res.RawJson, nil)
	urlFormat := "https://ip.ecnu.edu.cn/myip/format"
	res, err = HandleRequest(method, urlFormat, userAgent, "", timeout, nil)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, res.RawJson, nil)
}
