package rpc

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

var (
	HttpRequestErr       = errors.New("http请求异常")
	HttpStatusCodeBadErr = errors.New("返回值不为2xx 3xx")
)

type ValidResponseContent interface {
	Success() bool
}

func statusCodeCheck(response gorequest.Response) error {
	if response.StatusCode < http.StatusBadRequest {
		return nil
	}
	return HttpStatusCodeBadErr
}

func prepareRequest(base *gorequest.SuperAgent, headers map[string]string, params map[string]string) *gorequest.SuperAgent {
	for k, v := range headers {
		base.AppendHeader(k, v)
	}
	for k, v := range params {
		base = base.Param(k, v)
	}
	return base
}

func prepareRequestV2(base *gorequest.SuperAgent, headers map[string]string, params map[string]string, cookies []*http.Cookie) *gorequest.SuperAgent {
	if len(cookies) > 0 {
		base.AddCookies(cookies)
	}
	return prepareRequest(base, headers, params)
}

func handleRequestErrors(errs []error) error {
	if errs != nil {
		return HttpRequestErr
	}
	return nil
}

func unmarshal(body []byte, result interface{}) error {
	if err := json.Unmarshal(body, result); err != nil {
		return err
	}
	return nil
}

func processResponse(resp gorequest.Response, body []byte, errs []error, result interface{}, validators ...func(result interface{}) error) error {
	if err := handleRequestErrors(errs); err != nil {
		return err
	}

	statusBadError := statusCodeCheck(resp)

	if len(body) == 0 {
		// 204 NO CONTENT
		return nil
	}
	if err := unmarshal(body, &result); err != nil {
		return err
	}

	for _, fn := range validators {
		if err := fn(result); err != nil {
			return err
		}
	}
	return statusBadError
}

func Get(url string, headers map[string]string, params map[string]string, result interface{}, validators ...func(result interface{}) error) error {
	base := gorequest.New().Get(url)
	prepareRequest(base, headers, params)
	resp, body, errs := base.EndBytes()
	return processResponse(resp, body, errs, result, validators...)
}

func Post(url string, resource interface{}, headers map[string]string, params map[string]string, result interface{}, validators ...func(result interface{}) error) error {
	base := gorequest.New().Post(url)
	prepareRequest(base, headers, params)
	resp, body, errs := base.Send(resource).EndBytes()
	return processResponse(resp, body, errs, result, validators...)
}

func Delete(url string, headers map[string]string, params map[string]string, result interface{}, validators ...func(result interface{}) error) error {
	base := gorequest.New().Delete(url)
	prepareRequest(base, headers, params)
	resp, body, errs := base.EndBytes()
	return processResponse(resp, body, errs, result, validators...)
}

func doRequest(method string, url string, data interface{}, headers map[string]string, params map[string]string, cookies []*http.Cookie, result interface{}, validators ...func(result interface{}) error) (gorequest.Response, []byte, error) {
	base := gorequest.New().AddCookies(cookies)
	prepareRequestV2(base, headers, params, cookies)
	switch method {
	case "get":
		base = base.Get(url)
	case "post":
		base = base.Post(url)
	}
	resp, body, errs := base.Send(data).EndBytes()
	return resp, body, processResponse(resp, body, errs, result, validators...)
}

func GetV2(url string, headers map[string]string, params map[string]string, cookies []*http.Cookie, response interface{}, validators ...func(result interface{}) error) (*http.Response, []byte, error) {
	return doRequest("get", url, nil, headers, params, cookies, response, validators...)
}

func PostV2(url string, data interface{}, headers map[string]string, params map[string]string, cookies []*http.Cookie, response interface{}, validators ...func(result interface{}) error) (*http.Response, []byte, error) {
	return doRequest("post", url, data, headers, params, cookies, response, validators...)
}
