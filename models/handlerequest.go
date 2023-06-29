package models

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type DataReqMeta struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Proto  string `json:"proto"`
	Host   string `json:"host"`
}

type DataRespMeta struct {
	Proto  string `json:"proto"`
	Status string `json:"status"`
}

type DataResp struct {
	Meta   DataRespMeta      `json:"meta"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

type DataReq struct {
	Meta   DataReqMeta       `json:"meta"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

type Example struct {
	Curl string `json:"curl"`
}

type Resp struct {
	Request  DataReq     `json:"request"`
	Response DataResp    `json:"response"`
	RawJson  interface{} `json:"rawjson"`
	Example  Example     `json:"example"`
}

func HandleRequest(method, apiAddr, userAgent, reqBody string, timeout int, inputReqHeader map[string]string) (res Resp, err error) {
	req, resp, err := httpRequest(method, apiAddr, userAgent, timeout, bytes.NewBuffer([]byte(reqBody)), inputReqHeader)
	if err != nil {
		return
	}
	//处理返回body
	defer resp.Body.Close()
	//如何 resp.Header 中 Content-Encoding=gzip，直接读取会乱码，因此处理乱码
	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return
		}
	} else {
		reader = resp.Body
	}
	respBody, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	//处理Header
	respHeader := make(map[string]string)
	for k, v := range resp.Header {
		respHeader[k] = v[0]
	}
	reqHeader := make(map[string]string)
	for k, v := range req.Header {
		reqHeader[k] = v[0]
	}
	respProto := resp.Proto

	dataReqMeta := DataReqMeta{
		Method: req.Method,
		Path:   req.URL.RequestURI(),
		Proto:  req.Proto,
		Host:   req.Host,
	}
	dataReq := DataReq{
		Meta:   dataReqMeta,
		Header: reqHeader,
		Body:   base64.StdEncoding.EncodeToString([]byte(reqBody)),
	}
	dataRespMeta := DataRespMeta{
		Proto:  respProto,
		Status: resp.Status,
	}
	dataResp := DataResp{
		Meta:   dataRespMeta,
		Header: respHeader,
		Body:   base64.StdEncoding.EncodeToString(respBody),
	}
	res.Request = dataReq
	res.Response = dataResp

	var rawJson map[string]interface{}
	if err2 := json.Unmarshal([]byte(respBody), &rawJson); err2 == nil {
		res.RawJson = rawJson
	}

	res.Example.Curl = base64.StdEncoding.EncodeToString([]byte(convertRequestToCurl(req, reqBody)))
	return
}

func convertRequestToCurl(req *http.Request, body string) string {
	// 构建 curl 命令的基本部分
	curlCmd := fmt.Sprintf("curl -i -X %s", req.Method)

	// 添加请求头部
	for key, values := range req.Header {
		for _, value := range values {
			curlCmd += fmt.Sprintf(" -H '%s: %s'", key, value)
		}
	}

	// 添加请求 URL
	curlCmd += " '" + req.URL.String() + "'"

	// 添加请求体
	if body != "" {
		curlCmd += fmt.Sprintf(" -d '%s'", body)
	}

	return curlCmd
}
