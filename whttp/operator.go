package whttp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	noDefineMethod = "方法未定义"
)

type IHttpClient interface {
	Get(ctx context.Context, path, body string, header http.Header, timeout uint64,
		params map[string]string) (response *http.Response, err error)
	Post(ctx context.Context, path, body string, header http.Header, timeout uint64,
		params map[string]string) (response *http.Response, err error)
	Delete(ctx context.Context, path, body string, header http.Header, timeout uint64,
		params map[string]string) (response *http.Response, err error)
	Put(ctx context.Context, path, body string, header http.Header, timeout uint64,
		params map[string]string) (response *http.Response, err error)
	Patch(ctx context.Context, path, body string, header http.Header, timeout uint64,
		params map[string]string) (response *http.Response, err error)
	ResponseBody(ctx context.Context, method, path, body string, header http.Header, timeout uint64,
		params map[string]string) (int, []byte, error)
	Request(ctx context.Context, method, path, body string, header http.Header, timeout uint64,
		params map[string]string) (response *http.Response, err error)
	GetHttpClient(timeout uint64) *http.Client
}

type HttpClient struct {
}

func (h *HttpClient) Get(ctx context.Context, path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	return common(ctx, http.MethodGet, path, body, header, timeout, params)
}

func (h *HttpClient) Post(ctx context.Context, path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	return common(ctx, http.MethodPost, path, body, header, timeout, params)
}

func (h *HttpClient) Put(ctx context.Context, path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	return common(ctx, http.MethodPut, path, body, header, timeout, params)
}

func (h *HttpClient) Delete(ctx context.Context, path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	return common(ctx, http.MethodDelete, path, body, header, timeout, params)
}

func (h *HttpClient) Patch(ctx context.Context, path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	return common(ctx, http.MethodPatch, path, body, header, timeout, params)
}

func (h *HttpClient) Request(ctx context.Context, method, path, body string, header http.Header, timeout uint64,
	params map[string]string) (response *http.Response, err error) {
	switch method {
	case http.MethodGet:
		return h.Get(ctx, path, body, header, timeout, params)
	case http.MethodPut:
		return h.Put(ctx, path, body, header, timeout, params)
	case http.MethodPost:
		return h.Post(ctx, path, body, header, timeout, params)
	case http.MethodPatch:
		return h.Patch(ctx, path, body, header, timeout, params)
	case http.MethodDelete:
		return h.Delete(ctx, path, body, header, timeout, params)
	case http.MethodOptions, http.MethodTrace, http.MethodHead, http.MethodConnect:
		return common(ctx, method, path, body, header, timeout, params)
	default:
		return nil, fmt.Errorf("%s-%s", method, noDefineMethod)
	}
}

func (h *HttpClient) ResponseBody(ctx context.Context, method, path, body string, header http.Header, timeout uint64,
	params map[string]string) (int, []byte, error) {
	response, err := h.Request(ctx, method, path, body, header, timeout, params)
	if err != nil {
		return 0, nil, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return response.StatusCode, bytes, err
}

func (h *HttpClient) GetHttpClient(timeout uint64) *http.Client {
	return getClientConfig(timeout)
}
