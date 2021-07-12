package openapi

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	appId              string
	appSecret    	   string // 应用私钥
	apiDomain          string
	Client             *http.Client
	signName 		   string
	appidName          string
	dataName           string
	scheme             string
	timeout            time.Duration
	path               string
	header             *http.Header
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
func WithTimeout(timeout time.Duration) OptionFunc {
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

// 指定结果字段名称
func WithDataName(name string) OptionFunc{
	return func(c *Client) {
		c.dataName = name
	}
}

// 指定头部信息
func WithHeader(header *http.Header) OptionFunc {
	return func(c *Client) {
		c.header = header
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
	client.dataName = "response_data"
	client.timeout = 10 * time.Second
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// 获取URL请求参数
func (this *Client) URLValues(query *url.Values, params *Params) (*url.Values) {
	var p = &url.Values{}
	p.Add(this.appidName, this.appId)
	fmt.Println(UnixTime())
	p.Add("timestamp", fmt.Sprintf("%d", UnixTime()))
	if query != nil {
		for k, v := range *query{
			p.Add(k, v[0])
		}
	}
	s := params.Merge(p)
	sign := GenSign(this.appSecret, s)
	p.Add(this.signName, sign)
	return p
}

// 接口请求操作
func (this *Client) DoRequest(method string, api string, query *url.Values, data *Params, files *Params) (*Response, error) {
	var body io.Reader
	if data == nil {
		data = &Params{}
	}
	values := this.URLValues(query, data)
	header := make(http.Header)
	if this.header != nil {
		for k, v := range *this.header{
			header.Add(k, v[0])
		}
	}
	if files != nil {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		if data != nil {
			for k, v := range *data{
				writer.WriteField(k, v.String())
			}
		}
		for k, v := range *files{
			fn := v.String()
			st, err := os.Stat(fn)
			if err != nil {
				return nil, err
			}
			if st.IsDir() {
				return nil, errors.New("file is not eisted")
			}
			fp, _ := os.Open(fn)
			defer fp.Close()
			fs, _ := fp.Stat()
			fw, _ := writer.CreateFormFile(k, fs.Name())
			io.Copy(fw, fp)
		}
		header.Set("Content-Type", writer.FormDataContentType())
		writer.Close()
	}else{
		if data == nil || data.Size() == 0 {
			body = nil
		}else{
			body = data.Buffer()
			header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		}
	}
	if header.Get("User-Agent") == "" {
		header.Set("User-Agent", "HuaQiu OpenAPI Go SDK /0.1.0")
	}
	if this.Client == nil {
		tr := &http.Transport{
			TLSClientConfig:        &tls.Config{
				InsecureSkipVerify:          true,
			},
		}
		this.Client = &http.Client{Transport: tr, Timeout: this.timeout}
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
		rc, ok := body.(io.ReadCloser)
		if !ok && body != nil {
			rc = ioutil.NopCloser(body)
		}
		request.Body = rc
		switch v := body.(type) {
		case *bytes.Buffer:
			request.ContentLength = int64(v.Len())
			buf := v.Bytes()
			request.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
		default:

		}
	}

	request.Header = header
	request.URL.RawQuery = values.Encode()
	response, err := this.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, _ := ioutil.ReadAll(response.Body)
	resp := &Response{}
	resp.Response = *response
	resp.client = this
	resp.Content = content
	resp.Text = string(content)
	return resp, nil
}

// 设置UA
func (this *Client) UserAgent(useragent string) *Client {
	this.header.Set("User-Agent", useragent)
	return this;
}

// 设置超时时间
func (this *Client) Timeout(timeout time.Duration) *Client {
	this.timeout = timeout
	return this
}

// 设置头部信息
func (this *Client) Header(key, value string) *Client {
	this.header.Set(key, value)
	return this
}

// 接口 GET 请求
func (this *Client) Get(api string, params *Params) (*Response, error) {
	return this.DoRequest(http.MethodGet, api, params.UrlValues(), nil, nil)
}

// 接口 POST 请求
func (this *Client) Post(api string, params *Params)  (*Response, error) {
	return this.DoRequest(http.MethodPost, api, nil, params, nil)
}

// 接口 PUT 请求
func (this *Client) Put(api string, params *Params)  (*Response, error) {
	return this.DoRequest(http.MethodPut, api, nil, params, nil)
}

// 接口 DELETE 请求
func (this *Client) Delete(api string, params *Params)  (*Response, error) {
	return this.DoRequest(http.MethodDelete, api, params.UrlValues(), nil, nil)
}

// 接口 上传文件请求
func (this *Client) PostFile(api string, files *Params, params *Params) (*Response, error){
	if files == nil {
		return nil, errors.New("files is empty")
	}
	return this.DoRequest(http.MethodPost, api, nil, params, files)
}