package openapi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	appId              string
	appSecret    	   string // 应用私钥
	apiDomain          string
	Client             *http.Client
	signName 		   string
	appidName          string
	scheme             string
	timeout            int64
	path               string
}

type OptionFunc func(c *Client)

// 开启指定 http 客户端
func WithHTTPClient(client *http.Client) OptionFunc {
	return func(c *Client) {
		c.Client = client
	}
}

// 开启指定接口域名
func WithDomain(domain string) OptionFunc{
	return func(c *Client) {
		c.apiDomain = domain
	}
}

// 指定签名参数名称
func WithSignatureName(name string) OptionFunc {
	return func(c *Client) {
		c.signName = name
	}
}

// 指定协议
func WithScheme(scheme string) OptionFunc {
	return func(c *Client) {
		c.scheme = scheme
	}
}

// 指定超时时间
func WithTimeout(timeout int64) OptionFunc {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// 指定签名参数名称
func WithAppIdName(name string) OptionFunc {
	return func(c *Client) {
		c.appidName = name
	}
}


// 实例化 OpenAPI 客户端对象
func New(appid string, secret string, opts ... OptionFunc) (client *Client){
	client = &Client{}
	client.appId = appid
	client.appSecret = secret
	client.scheme = "http"
	client.signName  = "signature"
	client.appidName = "appid"
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// 获取URL请求参数
func (this *Client) URLValues(params Params) (url.Values) {
	var p = url.Values{}
	p.Add(this.appidName, this.appId)
	p.Add("timestamp", string(UnixTime()))
	s := params.Merge(p)
	sign := GenSign(this.appSecret, s)
	p.Add(this.signName, sign)
	return p
}

// 接口请求操作
func (this *Client) DoRequest(method string, api string, query url.Values, data Params, header http.Header) (*Response, error) {
	var body io.ReadCloser
	if data == nil {
		data = Params{}
		body = nil
	}else{
		body = data.Body()
	}
	values := this.URLValues(data)
	if query == nil {
		query = make(url.Values)
	}
	for k, v := range values {
		query[k] = v
	}
	if header == nil {
		header = make(http.Header)
	}
	if this.Client == nil {
		this.Client = &http.Client{Timeout: time.Duration(this.timeout) * time.Second}
	}
	request := &http.Request{
		Method: method,
		URL: &url.URL{
			Scheme: this.scheme,
			Host:   this.apiDomain,
			Path:   fmt.Sprintf("%s%s", this.path, api),
		},
		Proto:      "HTTP/1.1", // 29
		ProtoMajor: 1,
		ProtoMinor: 1,
		Host:       this.apiDomain,
	}
	if body != nil {
		request.Body = body
	}
	request.URL.RawQuery = query.Encode()
	response, err := this.Client.Do(request)
	if err != nil {
		return nil, err
	}
	resp := &Response{}
	resp.Response = *response;
	return resp, nil
}

// 接口 GET 请求
func (this *Client) Get(api string, params Params) (*Response, error) {
	return this.DoRequest(http.MethodGet, api, params.UrlValues(), nil, nil)
}

// 接口 POST 请求
func (this *Client) Post(api string, params Params)  (*Response, error) {
	return this.DoRequest(http.MethodPost, api, nil, params, nil)
}

// 接口 PUT 请求
func (this *Client) Put(api string, params Params)  (*Response, error) {
	return this.DoRequest(http.MethodPut, api, nil, params, nil)
}

// 接口 DELETE 请求
func (this *Client) Delete(api string, params Params)  (*Response, error) {
	return this.DoRequest(http.MethodDelete, api, params.UrlValues(), nil, nil)
}