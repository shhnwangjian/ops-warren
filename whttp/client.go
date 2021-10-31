package whttp

import (
	"bytes"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	maxIdleConns        = 100
	maxIdleConnsPerHost = 100
	idleConnTimeout     = 30
	keepAlive           = 30
	netTimeout          = 30
)

func common(method, path, body string, header http.Header, timeout uint64, params map[string]string) (response *http.Response, err error) {
	client := getClientConfig(timeout)
	req, err := http.NewRequest(method, buildUrl(path, params), strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	response = resp
	return
}

func getClientConfig(timeout uint64) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: false, // 是否开启http keepalive功能，也即是否重用连接，默认开启(false)
			Proxy:             http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   netTimeout * time.Second, // 限制建立TCP连接所花费的时间
				KeepAlive: keepAlive * time.Second,
			}).DialContext,
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			IdleConnTimeout:     idleConnTimeout * time.Second, // 空闲timeout设置，也即socket在该时间内没有交互则自动关闭连接
		},
		Timeout: time.Duration(timeout) * time.Second, // 客户端发出的请求的时间限制。该超时包括连接时间、任何、重定向，以及读取响应体
	}
}

func buildUrl(path string, params map[string]string) string {
	if len(params) == 0 {
		return path
	}
	if !strings.HasSuffix(path, "?") {
		path = path + "?"
	}
	var buf bytes.Buffer
	for key, value := range params {
		buf.WriteString(url.QueryEscape(key))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(value))
		buf.WriteByte('&')
	}
	paramBody := buf.String()
	if strings.HasSuffix(paramBody, "&") {
		paramBody = paramBody[:len(paramBody)-1]
	}
	return path + paramBody
}
