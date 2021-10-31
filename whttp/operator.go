package whttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	noDefineMethod = "方法未定义"
)

type IHttpClient interface {
	Get(path, body string, header http.Header, timeout uint64, params map[string]string) (response *http.Response, err error)
	Post(path, body string, header http.Header, timeout uint64, params map[string]string) (response *http.Response, err error)
	Delete(path, body string, header http.Header, timeout uint64, params map[string]string) (response *http.Response, err error)
	Put(path, body string, header http.Header, timeout uint64, params map[string]string) (response *http.Response, err error)
	Patch(path, body string, header http.Header, timeout uint64, params map[string]string) (response *http.Response, err error)
	ResponseBody(method, path, body string, header http.Header, timeout uint64, params map[string]string) (int, []byte, error)
	Request(method, path, body string, header http.Header, timeout uint64, params map[string]string) (response *http.Response, err error)
}

type HttpClient struct {
}

func (h *HttpClient) Get(path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	return common(http.MethodGet, path, body, header, timeout, params)
}

func (h *HttpClient) Post(path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	return common(http.MethodPost, path, body, header, timeout, params)
}

func (h *HttpClient) Put(path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	return common(http.MethodPut, path, body, header, timeout, params)
}

func (h *HttpClient) Delete(path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	return common(http.MethodDelete, path, body, header, timeout, params)
}

func (h *HttpClient) Patch(path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	return common(http.MethodPatch, path, body, header, timeout, params)
}

func (h *HttpClient) Request(method, path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	switch method {
	case http.MethodGet:
		return h.Get(path, body, header, timeout, params)
	case http.MethodPut:
		return h.Put(path, body, header, timeout, params)
	case http.MethodPost:
		return h.Post(path, body, header, timeout, params)
	case http.MethodPatch:
		return h.Patch(path, body, header, timeout, params)
	case http.MethodDelete:
		return h.Delete(path, body, header, timeout, params)
	case http.MethodOptions, http.MethodTrace, http.MethodHead, http.MethodConnect:
		return common(method, path, body, header, timeout, params)
	default:
		return nil, fmt.Errorf("%s-%s", method, noDefineMethod)
	}
}

func (h *HttpClient) ResponseBody(method, path, body string, header http.Header, timeout uint64,
	params map[string]string) (int, []byte, error) {
	response, err := h.Request(method, path, body, header, timeout, params)
	if err != nil {
		return 0, nil, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return response.StatusCode, bytes, err
}
