package models

import (
	"crypto/tls"
	"io"
	"time"

	"net/http"
)

func httpRequest(method, url, userAgent string, timeout int, data io.Reader, headers map[string]string) (req *http.Request, resp *http.Response, err error) {
	req, err = http.NewRequest(method, url, data)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Close = true
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(timeout),
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	return
}
